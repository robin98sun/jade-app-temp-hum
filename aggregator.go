package main

import (
	// "log"
	// "time"
	// "uta.edu/aces/jadesdk"
)

type Aggregator struct {
}

func NewAggregator() *Aggregator {
	return &Aggregator{}
}

type AggregatorInput struct {
	Amount int `json:"amount,omitempty"`
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
	subtaskResult := subtaskResultInst.(*AggregatorInput)
	result := &AggregatorInput{
		Amount: subtaskResult.Amount,
	}
	var cumulation *AggregatorInput
	if cumulationInst != nil {
		cumulation = cumulationInst.(*AggregatorInput)
	}
	// aggregate subtasks
	if cumulation != nil {
		result.Amount += cumulation.Amount
	}

	return result, nil
}
