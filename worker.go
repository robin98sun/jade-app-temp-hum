package main

import (
	// "math/rand"
	// "sort"
	"time"
	"log"
	"uta.edu/aces/jadesdk"
	"net/http"
	"encoding/json"
	"strings"
	"bytes"
)

type Worker struct {
	SDK *jadesdk.JadeSDK
}

func NewWorker(sdk *jadesdk.JadeSDK) *Worker {
	return &Worker{
		SDK: sdk,
	}
}

type WorkerInput struct {
	Days  int `json:"days,omitempty"`
	StartDate string `json:"startDate,omitempty"`
	EndDate string `json:"endDate,omitempty"`
}

func (w *Worker) ShapeInput() interface{} {
	return &WorkerInput{}
}

/* input example
{
   days: 10,
   startDate: '2021-05-11',
   endDate: '2021-05-12'
}
*/

type Request struct {
	StartDate string `json:"date3,omitempty"`
	EndDate string `json:"date4,omitempty"`
	FetchTemperature string `json:"temp_box,omitempty"`
	FetchHumidity string `json:"hum_box,omitempty"`
}

type Response struct {
	Time string `json:"TIME,omitempty"`
	Temperature float64 `json:"Temperature,omitempty"`
	Humidity float64 `json:"Humidity,omitempty"`
}

// the input is WorkerInput, output is AggregatorInput
func (w *Worker) Handler(inputInst interface{}) (interface{}, error) {
	var input *WorkerInput
	input = inputInst.(*WorkerInput)

	log.Printf("input: startDate: %v, endDate: %v, days: %v", input.StartDate, input.EndDate, input.Days)

	startDate := input.StartDate
	endDate := input.EndDate
	days := input.Days

	if days < 1 {
		days = 1
	}

	endDateTime := time.Now()
	startDateTime := time.Now()
	if startDate == "" && endDate == "" {
		startDateTime = endDateTime.AddDate(0,0, -days)
	} else if endDate == "" {
		startDateTime, _ = time.Parse("1999-01-01", startDate)
		if input.Days > 0 {
			endDateTime = startDateTime.AddDate(0,0, days)
		}
	} else if startDate == "" {
		endDateTime, _ = time.Parse("1999-01-01", endDate)
		startDateTime = endDateTime.AddDate(0,0, -days)
	} else {
		startDateTime, _ = time.Parse("1999-01-01", startDate)
		endDateTime, _ = time.Parse("1999-01-01", endDate)
	}

	startDate = startDateTime.Format("1999-01-01")
	endDate = endDateTime.Format("1999-01-01")

	// do some job
	capaName := "jade-app-temp-hum"
	var capability *jadesdk.Capability
	if w.SDK != nil && w.SDK.Conf.Capabilities != nil && len(w.SDK.Conf.Capabilities) > 0 {
		for _, cap := range w.SDK.Conf.Capabilities {
			if cap.Name == capaName {
				capability = cap
			}
		}
	}

	fetchedData := []*Response{}
	if capability != nil {
		action := capability.Action
		url := capability.URL
		log.Printf("action: %v, url: %v", action, url)
		// Send the register information to upper node
		payload := &Request{
			StartDate: startDate,
			EndDate: endDate,
			FetchTemperature: "temp",
			FetchHumidity: "hum",
		}
		reqbody, _ := json.Marshal(payload)
		req, err := http.NewRequest(strings.ToUpper(action), url, bytes.NewBuffer(reqbody))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		client := &http.Client{}
		res, err := client.Do(req)
		if err == nil && res.Body != nil{
			resData := []*Response{}
			json.NewDecoder(res.Body).Decode(&resData)
			fetchedData = resData
			log.Printf("SUCCESSFULLY fetched data amount: %v", len(fetchedData))
		} else if res.Body == nil {
			log.Println("ERROR: response does not have a body")
		} else {
			log.Printf("ERROR: error when requesting web service: %v \n", err)
		}
	}

	forwardToAggregator := &AggregatorInput{
		Amount: len(fetchedData),
	}

	// done
	return forwardToAggregator, nil
}
