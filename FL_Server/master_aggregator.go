package main

import (
	serverMsg "fl_system/messages"
	"log"
	"strconv"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type MasterAggregator struct {
	AggregatorFinishedTraining int
}

func (m *MasterAggregator) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		log.Printf("[+] MASTER AGGREGATOR ::::: Master Aggregator has started\n")

	case *actor.Stopping:
		log.Printf("[+] MASTER AGGREGATOR ::::: Stopping Master Aggregator, is about to shut down\n")

	case *actor.Restarting:
		log.Printf("[+]MASTER AGGREGATOR ::::: Restarting Master Aggregator\n")

	case *actor.Stopped:
		log.Printf("[+]MASTER AGGREGATOR ::::: Master Aggregator and all its children are stopped\n")

	case *serverMsg.SpawnNewAggr:
		log.Printf("[+]MASTER AGGREGATOR ::::: Master Aggregator Received Command to Spawn Aggregator : NumAggr : %+v\n", msg.NumAggr)

		for i := 0; i < msg.NumAggr; i++ {
			aggrProps := actor.PropsFromProducer(newAggregator)
			time.Sleep(2 * time.Second)
			ctx.SpawnNamed(aggrProps, "Aggr"+strconv.Itoa(i))
		}
		// log.Print(ctx.Children())
		time.Sleep(1 * time.Second)
		log.Printf("[+]MASTER AGGREGATOR ::::: Aggregator Spawn Successful\n")
		ctx.Respond(&serverMsg.AggrSpawnSuccessful{AggrPidList: ctx.Children()})
		time.Sleep(1 * time.Second)

	case *serverMsg.ModelTrainingCompletion:
		m.AggregatorFinishedTraining += 1
		if m.AggregatorFinishedTraining == len(ctx.Children()) {
			for i := 0; i < len(ctx.Children()); i++ {
				time.Sleep(2 * time.Second)
				ctx.Stop(ctx.Children()[i])
			}
			ctx.Request(ctx.Parent(), &serverMsg.ModelTrainingCompletion{})
		}

	default:
	}
}

func newMasterAggregator() actor.Actor {
	return &MasterAggregator{}
}
