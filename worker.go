package main

import (
	// "math/rand"
	// "sort"
	// "time"
	"log"
	"uta.edu/aces/jadesdk"
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

// the input is WorkerInput, output is AggregatorInput
func (w *Worker) Handler(inputInst interface{}) (interface{}, error) {
	var input *WorkerInput
	input = inputInst.(*WorkerInput)

	log.Printf("input: startDate: %v, endDate: %v, days: %v", input.StartDate, input.EndDate, input.Days)

	capaName := "jade-app-temp-hum"
	var capability *jadesdk.Capability
	if w.SDK != nil && w.SDK.Conf.Capabilities != nil && len(w.SDK.Conf.Capabilities) > 0 {
		for _, cap := range w.SDK.Conf.Capabilities {
			if cap.Name == capaName {
				capability = cap
			}
		}
	}
	if capability != nil {
		action := capability.Action
		url := capability.URL
		log.Printf("action: %v, url: %v", action, url)
	}

	forwardToAggregator := &AggregatorInput{
	}
	// do some job
	

	// done
	return forwardToAggregator, nil
}
