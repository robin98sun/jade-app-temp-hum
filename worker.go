package main

import (
	"math/rand"
	"sort"
	"time"
)

type Worker struct {
}

func NewWorker() *Worker {
	return &Worker{}
}

type WorkerInput struct {
	Days  int `json:"days,omitempty"`
	StartDate string `json:"startDate,omitempty"`
	EndDate string `json:"endDate,omitempty"`
}

type AggregatorInput struct {
	Cmd    string  `json:"cmd,omitempty"`
	EatTime   int64  `json:"eatTime,omitempty"`
	Size      int64  `json:"size,omitempty"`
	Pieces []int64 `json:"pieces,omitempty"`
	DigestTime int64 `json:"digestTime,omitempty"`
	DigestFactor int64 `json:"digestFactor,omitempty"`
}

func (w *Worker) ShapeInput() interface{} {
	return &WorkerInput{}
}

// the input is WorkerInput, output is AggregatorInput
func (w *Worker) Handler(inputInst interface{}) (interface{}, error) {
	var input *WorkerInput
	input = inputInst.(*WorkerInput)

	feedStomach := &AggregatorInput{
		Cmd:    input.Cmd,
		EatTime:   input.Size,
	}
	// do some job
	if (input.Cmd == "gen and merge" || input.Cmd == "gen and merge and wait") && input.Size > 0 {
		rand.Seed(time.Now().UTC().UnixNano())
		pieces := []int64{}
		for i := int64(0); i < input.Size; i++ {
			n := rand.Int63n(input.Size * 100)
			pieces = append(pieces, n)
		}
		sort.Slice(pieces, func(i, j int) bool {
			return pieces[i] < pieces[j]
		})
		feedStomach.Pieces = pieces
		feedStomach.EatTime = 0
	} 

	// wait some time
	if input.Cmd == "service time" || input.Cmd == "gen and merge and wait" {
		rand.Seed(time.Now().UTC().UnixNano())
		n := int64(0)
		if input.MaxEatTime >= input.MinEatTime && input.MinEatTime > 0 {
			if input.MinEatTime == input.MaxEatTime {
				n = input.MinEatTime
			} else {
				n = rand.Int63n(input.MaxEatTime-input.MinEatTime)
				n += int64(input.MinEatTime)
			}
		}
		time.Sleep(time.Duration(n) * time.Millisecond)
		feedStomach.EatTime = n	
	}

	// forward digest options
	feedStomach.DigestTime = input.DigestTime
	feedStomach.DigestFactor = input.DigestFactor

	// done
	return feedStomach, nil
}
