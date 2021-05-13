package main

import (
	"uta.edu/aces/jadesdk"
)

func main() {
	jade := jadesdk.NewJadeSDK()
	jade.Verbose(false)

	workerMod := NewWorker(jade.Conf.Capabilities)
	jade.SetDefaultWorkerModule(workerMod)

	aggregatorMod := NewAggregator(jade.Conf.Capabilities)
	jade.SetDefaultAggregatorModule(aggregatorMod)
	jade.CreateHTTPServer()
}
