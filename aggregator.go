package main

import (
	// "log"
	// "time"
	"uta.edu/aces/jadesdk"
)

type Aggregator struct {
}

func NewAggregator(capabilities []*jadesdk.Capability) *Aggregator {
	return &Aggregator{}
}

type AggregatorInput struct {
	Data []interface{} `json:"data,omitempty"`
}

func (w *Aggregator) ShapeResultOfSubtask() interface{} {
	return &AggregatorInput{}
}

func (w *Aggregator) ShapeCumulation() interface{} {
	return &AggregatorInput{}
}

func (w *Aggregator) Handler(cumulationInst interface{}, previousResults []interface{}, subtaskResultInst interface{}) (interface{}, error) {
	if subtaskResultInst == nil {
		return cumulationInst, nil
	}
	// subtaskResult := subtaskResultInst.(*AggregatorInput)
	result := &AggregatorInput{
	}
	// var cumulation *AggregatorInput
	// if cumulationInst != nil {
	// 	cumulation = cumulationInst.(*AggregatorInput)
	// }

	// aggregate subtasks
	

	return result, nil
}
