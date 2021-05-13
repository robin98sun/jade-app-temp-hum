package main

import (
	"uta.edu/aces/jadesdk"
)

func main() {
	jade := jadesdk.NewJadeSDK()
	jade.Verbose(true)

	workerMod := NewWorker(jade.Conf.Capabilities)
	jade.SetDefaultWorkerModule(workerMod)

	aggregatorMod := NewAggregator()
	jade.SetDefaultAggregatorModule(aggregatorMod)
	jade.CreateHTTPServer()
}
