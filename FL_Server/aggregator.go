package main

import (
	serverMsg "fl_system/messages"
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type Aggregator struct {
	TrainersUnderAggr         []*actor.PID
	TrainersCompletedTraining int
}

func (a *Aggregator) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		log.Printf("[+]AGGREGATOR ::::::: Aggregator %v has started\n", ctx.Self().Id)

	case *actor.Stopping:
		log.Printf("[+]AGGREGATOR ::::::: Aggregator, is about to shut down\n")

	case *actor.Restarting:
		log.Printf("[+]AGGREGATOR ::::::: Restarting Aggregator\n")

	case *actor.Stopped:
		log.Printf("[+]AGGREGATOR ::::::: Aggregator and all its children are stopped\n")

	case *serverMsg.ReceiveDevices:
		log.Printf("AGGREGATOR ::::::: Receiving devices from Selector\n")
		log.Printf("\n AGGREGATOR ::::::: AggregatorPID: %v\n", ctx.Self().Id)
		log.Printf("AGGREGATOR ::::::: Devices Received: %v", msg.DevicesPidList)
		a.TrainersUnderAggr = msg.DevicesPidList
		time.Sleep(1 * time.Second)

		ctx.Request(msg.DevicesPidList[0], &serverMsg.Demo{})

	case *serverMsg.ModelCarryingMessage:
		log.Printf("AGGREGATOR ::::::: Model Received........Sending models to trainer")
		log.Print(len(a.TrainersUnderAggr))
		for i := 0; i < len(a.TrainersUnderAggr); i++ {
			time.Sleep(1 * time.Second)
			ctx.Request(a.TrainersUnderAggr[i], &serverMsg.ModelCarryingMessage{Model: msg.Model})
		}

	case *serverMsg.ModelTrainingCompletion:
		a.TrainersCompletedTraining += 1
		if a.TrainersCompletedTraining == len(a.TrainersUnderAggr) {
			ctx.Request(ctx.Parent(), &serverMsg.ModelTrainingCompletion{})
		}

	default:
	}
}

func newAggregator() actor.Actor {
	return &Aggregator{}
}
