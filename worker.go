package main

import (
	// "math/rand"
	// "sort"
	// "time"
	"log"
	"uta.edu/aces/jadesdk"
)

type Worker struct {
	Capabilities []*jadesdk.Capability
}

func NewWorker(capabilities []*jadesdk.Capability) *Worker {
	return &Worker{
		Capabilities: capabilities,
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

// the input is WorkerInput, output is AggregatorInput
func (w *Worker) Handler(inputInst interface{}) (interface{}, error) {
	var input *WorkerInput
	input = inputInst.(*WorkerInput)

	log.Printf("input: startDate: %v, endDate: %v, days: %v", input.StartDate, input.EndDate, input.Days)

	for i, cap := range w.Capabilities {
		log.Print("capability[%v] name: %v, value: %v, api: %v, type: %v, action: %v, url: %v", 
			i, cap.Name, cap.Value, cap.API, cap.Type, cap.Action, cap.URL,
		)
		if cap.Parameters != nil && len(cap.Parameters) > 0 {
			for j, param := range cap.Parameters {
				log.Print("   param[%v] name: %v, type: %v", param.Name, param.Type)
			}
		}
	}


	forwardToAggregator := &AggregatorInput{
	}
	// do some job
	

	// done
	return forwardToAggregator, nil
}
