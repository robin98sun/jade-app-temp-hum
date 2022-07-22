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
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

type Worker struct {
	SDK *jadesdk.JadeSDK
	DbConn *sql.DB
}

func NewWorker(sdk *jadesdk.JadeSDK) *Worker {
	return &Worker{
		SDK: sdk,
	}
}

type TargetType string

const (
	TargetTypeDatabase TargetType = "database"
	TargetTypeService             = "service"
)

type WorkerInput struct {
	Days  int `json:"days,omitempty"`
	StartDate string `json:"startDate,omitempty"`
	EndDate string `json:"endDate,omitempty"`
	Target TargetType `json:"target,omitempty"`
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
	Preprocessing float64 `json:"preprocessing,omitempty"`
	Connection float64 `json:"connection,omitempty"`
	Query float64 `json:"query,omitempty"`
	Results []*DataItem
	DBHost string `json:"dbhost,omitempty"`
}


// the input is WorkerInput, output is AggregatorInput
func (w *Worker) Handler(inputInst interface{}) (interface{}, error) {
	startTime := time.Now()
	var input *WorkerInput
	input = inputInst.(*WorkerInput)

	startDate := input.StartDate
	endDate := input.EndDate
	days := input.Days
	target := input.Target
	if target == "" {
		target = TargetTypeDatabase
	}

	// if days < 1 {
	// 	days = 1
	// }

	endDateTime := time.Now()
	startDateTime := time.Now()
	// formatStr := "2006-01-02"
	formatStr := "2006-01-02T15:04:05Z07:00"
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
	var capability_service *jadesdk.Capability
	db_name := ""
	db_pass := ""
	db_user := ""
	db_host := ""
	db_port := ""
	if w.SDK != nil && w.SDK.Conf.Capabilities != nil && len(w.SDK.Conf.Capabilities) > 0 {
		for _, cap := range w.SDK.Conf.Capabilities {
			if cap.Name == "jade-app-temp-hum_service" {
				capability_service = cap
			} else if cap.Name == "jade-app-temp-hum_db_host" {
				db_host = cap.Value
			} else if cap.Name == "jade-app-temp-hum_db_user" {
				db_user = cap.Value
			} else if cap.Name == "jade-app-temp-hum_db_name" {
				db_name = cap.Value
			} else if cap.Name == "jade-app-temp-hum_db_pass" {
				db_pass = cap.Value
			} else if cap.Name == "jade-app-temp-hum_db_port" {
				db_port = cap.Value
			}
		}
	}

	fetchedData := &Response{}
	conn_desc := "N/A"
	dbhost := ""
	if target == TargetTypeDatabase {
		if db_name != "" && db_pass != "" && db_user != "" && db_host != "" && db_port != "" {
			db_desc := db_user + ":<pass>@tcp("+ db_host + ":" + db_port + ")/" + db_name
			db_conn_str := db_user + ":" + db_pass + "@tcp("+ db_host + ":" + db_port + ")/" + db_name
			conn_desc = db_desc
			startConnTime := time.Now()
			// connect database
			if w.DbConn == nil {
				db, err := sql.Open("mysql", db_conn_str)
				if err != nil {
					fetchedData.Error = "connection error: " + err.Error()
				} else {
					w.DbConn = db
				}
			}
			startQueryTime := time.Now()
			if w.DbConn != nil && fetchedData.Error == "" {
				results, err := w.DbConn.Query("SELECT TIME, VALUE_2 AS 'Temperature' , VALUE_1 AS 'Humidity' FROM sensor WHERE TIME BETWEEN ? AND ?", startDate, endDate)
				if err != nil {
					fetchedData.Error = "query error: " + err.Error()
				} else {
					for results.Next() {
						var item DataItem
						err = results.Scan(&item.Time, &item.Temperature, &item.Humidity)
						if err != nil {
							fetchedData.Error += " [data error]: " + err.Error() + "; "
						} else {
							if fetchedData.Results == nil {
								fetchedData.Results = []*DataItem{}
							}
							fetchedData.Results = append(fetchedData.Results, &item)
						}
					}
				}
			}
			fetchedData.Query = float64(time.Now().Sub(startQueryTime))/float64(time.Millisecond)
			fetchedData.Connection = float64(startQueryTime.Sub(startConnTime))/float64(time.Millisecond)
			fetchedData.Preprocessing = float64(startConnTime.Sub(startTime))/float64(time.Millisecond)
		} else {
			log.Printf("ERROR: capabilities for [%v] not found or incomplete", target)
		}
		
	} else if target == TargetTypeService {
		if capability_service != nil {
			action := capability_service.Action
			serviceUrl := capability_service.URL
			conn_desc = serviceUrl

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
			if err == nil && res != nil && res.Body != nil{
				resData := &Response{}
				bodyDecoder := json.NewDecoder(res.Body)
				if bodyDecoder != nil {
					bodyDecoder.Decode(&resData)
				} else {
					log.Printf("ERROR: can not decode body from response, {%v}", serviceUrl)
				}
				fetchedData = resData
				dbhost = "@" + resData.DBHost
				// log.Printf("SUCCESSFULLY fetched data amount: %v", len(fetchedData))
			} else if res == nil {
				log.Printf("ERROR: response is null, {%v}", serviceUrl)
			} else if res.Body == nil {
				log.Printf("ERROR: response does not have a body, {%v}", serviceUrl)
			} else {
				log.Printf("ERROR: error when requesting web service {%v}: %v \n", serviceUrl, err)
			}
		} else {
			log.Printf("ERROR: capability for [%v] not found", target)
		}
	}
	

	forwardToAggregator := &AggregatorInput{
		Amount: len(fetchedData.Results),
		AvgTemp: 0,
		AvgHum: 0,
	}

	if len(fetchedData.Results) > 0 {
		for _, item := range fetchedData.Results {
			forwardToAggregator.AvgTemp += item.Temperature
			forwardToAggregator.AvgHum += item.Humidity
		}
		forwardToAggregator.AvgTemp /= float64(len(fetchedData.Results))
		forwardToAggregator.AvgHum /= float64(len(fetchedData.Results))
	}

	log.Printf("startDate: %v, endDate: %v, days: %v, fetched lines: %v, preprocessing time(ms): %v, connection time (ms): %v, query time (ms): %v, [%v%v]:{%v}", 
		startDate, endDate, days, 
		len(fetchedData.Results),
		fetchedData.Preprocessing,
		fetchedData.Connection,
		fetchedData.Query,
		target, dbhost, conn_desc,
	)
	if fetchedData.Error != "" {
		log.Printf("%v ERROR: startDate: %v, endDate: %v, days: %v, fetched lines: %v, [%v%v]:{%v}, ERROR: %v", 
			target,
			startDate, endDate, days, 
			len(fetchedData.Results),
			target, dbhost, conn_desc,
			fetchedData.Error,
		)
	}

	// done
	return forwardToAggregator, nil
}
