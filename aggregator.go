package main

import (
	"log"
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
	AvgTemp float64 `json:"avg_temp,omitempty"`
	AvgHum float64 `json:"avg_hum,omitempty"`
}

func (a *AggregatorInput) Copy() *AggregatorInput {
	if a == nil {return nil}
	return &AggregatorInput{
		Amount: a.Amount,
		AvgTemp: a.AvgTemp,
		AvgHum: a.AvgHum,
	}
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
	result := subtaskResult.Copy()

	var cumulation *AggregatorInput
	if cumulationInst != nil {
		cumulation = cumulationInst.(*AggregatorInput)
	}
	// aggregate subtasks
	if cumulation != nil && result != nil {
		if result.Amount + cumulation.Amount > 0 {
			result.AvgTemp = cumulation.AvgTemp * float64(cumulation.Amount) + result.AvgTemp * float64(result.Amount)
			result.AvgHum = cumulation.AvgHum * float64(cumulation.Amount) + result.AvgHum * float64(result.Amount)
			result.AvgTemp /= float64(result.Amount + cumulation.Amount)
			result.AvgHum /= float64(result.Amount + cumulation.Amount)
		}
		result.Amount += cumulation.Amount
	}

	log.Printf("task result: %v, cumulative result: %v",
		subtaskResult.Amount, result.Amount,
	)

	return result, nil
}
