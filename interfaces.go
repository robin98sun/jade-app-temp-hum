package main

import (
	"uta.edu/aces/jadesdk"
)

func main() {
	jade := jadesdk.NewJadeSDK()
	jade.Verbose(false)

	workerMod := NewWorker()
	jade.SetDefaultWorkerModule(workerMod)

	aggregatorMod := NewAggregator()
	jade.SetDefaultAggregatorModule(aggregatorMod)
	jade.CreateHTTPServer()
}
