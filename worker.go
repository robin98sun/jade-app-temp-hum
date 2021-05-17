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
    "net/url"
    "strconv"
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

type DataItem struct {
	Time string `json:"TIME,omitempty"`
	Temperature float64 `json:"Temperature,omitempty"`
	Humidity float64 `json:"Humidity,omitempty"`
}

type Response struct {
	Error string `json:"error,omitempty"`
	Preprocessing int64 `json:"preprocessing,omitempty"`
	Connection int64 `json:"connection,omitempty"`
	Query int64 `json:"query,omitempty"`
	Results []*DataItem
}

// the input is WorkerInput, output is AggregatorInput
func (w *Worker) Handler(inputInst interface{}) (interface{}, error) {
	var input *WorkerInput
	input = inputInst.(*WorkerInput)

	startDate := input.StartDate
	endDate := input.EndDate
	days := input.Days

	if days < 1 {
		days = 1
	}

	endDateTime := time.Now()
	startDateTime := time.Now()
	formatStr := "2006-01-02"
	if startDate == "" && endDate == "" {
		startDateTime = endDateTime.AddDate(0,0, -days)
	} else if endDate == "" {
		startDateTime, _ = time.Parse(formatStr, startDate)
		if input.Days > 0 {
			endDateTime = startDateTime.AddDate(0,0, days)
		}
	} else if startDate == "" {
		endDateTime, _ = time.Parse(formatStr, endDate)
		startDateTime = endDateTime.AddDate(0,0, -days)
	} else {
		if startDateTime.Sub(endDateTime) > 0 {
			endDateTime = startDateTime.AddDate(0,0, days)
		}
		startDateTime, _ = time.Parse(formatStr, startDate)
		endDateTime, _ = time.Parse(formatStr, endDate)
	}

	startDate = startDateTime.Format(formatStr)
	endDate = endDateTime.Format(formatStr)

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

	fetchedData := &Response{}
	serviceUrl := "N/A"
	if capability != nil {
		action := capability.Action
		serviceUrl = capability.URL
		// log.Printf("action: %v, url: %v, startDate: %v, endDate: %v, days: %v", action, serviceUrl, startDate, endDate, days)
		// Send the register information to upper node
		payload := url.Values{}
		payload.Set("date3", startDate)
		payload.Set("date4", endDate)
		payload.Set("temp_box", "temp")
		payload.Set("hum_box", "hum")

		req, err := http.NewRequest(strings.ToUpper(action), serviceUrl, strings.NewReader(payload.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(payload.Encode())))

		client := &http.Client{}
		res, err := client.Do(req)
		if err == nil && res.Body != nil{
			resData := &Response{}
			bodyDecoder := json.NewDecoder(res.Body)
			if bodyDecoder != nil {
				bodyDecoder.Decode(&resData)
			} else {
				log.Printf("ERROR: can not decode body from response, {%v}", serviceUrl)
			}
			fetchedData = resData
			// log.Printf("SUCCESSFULLY fetched data amount: %v", len(fetchedData))
		} else if res.Body == nil {
			log.Printf("ERROR: response does not have a body, {%v}", serviceUrl)
		} else {
			log.Printf("ERROR: error when requesting web service {%v}: %v \n", serviceUrl, err)
		}
	}

	forwardToAggregator := &AggregatorInput{
		Amount: len(fetchedData.Results),
	}
	log.Printf("startDate: %v, endDate: %v, days: %v, fetched lines: %v, preprocessing time(ms): %v, connection time (ms): %v, query time (ms): %v, {%v}", 
		input.StartDate, input.EndDate, input.Days, 
		len(fetchedData.Results),
		fetchedData.Preprocessing,
		fetchedData.Connection,
		fetchedData.Query,
		serviceUrl,
	)
	if fetchedData.Error != "" {
		log.Printf("SERVICE ERROR: startDate: %v, endDate: %v, days: %v, fetched lines: %v, {%v}, ERROR: %v", 
			input.StartDate, input.EndDate, input.Days, 
			len(fetchedData.Results),
			serviceUrl,
			fetchedData.Error,
		)
	}

	// done
	return forwardToAggregator, nil
}
