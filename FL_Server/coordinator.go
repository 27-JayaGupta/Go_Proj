package main

import (
	serverMsg "fl_system/messages"
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type coordinatorActor struct {
	devices_needed        int
	devices_for_each_aggr int
	selectorPid           *actor.PID
	AggrPidList           []*actor.PID
}

func (c *coordinatorActor) Receive(ctx actor.Context) {
	var numAggr int = c.devices_needed / c.devices_for_each_aggr
	switch msg := ctx.Message().(type) {

	case *serverMsg.SelectorInfoMsg:
		log.Printf("[+]Received Info from selector : devices_avaialablle : %+v\n", msg.Devices_available)
		c.selectorPid = ctx.Sender()
		if msg.Devices_available < c.devices_needed {
			time.Sleep(1 * time.Second)
			ctx.Respond(&serverMsg.NotEnoughDevicesAvailable{})
			return
		}

		if test_mode == true {
			return
		}

		masterAggregatorProps := actor.PropsFromProducer(newMasterAggregator)
		masterAggregatorPid, err := ctx.SpawnNamed(masterAggregatorProps, "master_aggregator")
		if err != nil {
			log.Print(err.Error())
			return
		}
		time.Sleep(1 * time.Second)

		ctx.Request(masterAggregatorPid, &serverMsg.SpawnNewAggr{NumAggr: numAggr})
		time.Sleep(1 * time.Second)

	case *serverMsg.AggrSpawnSuccessful:
		time.Sleep(1 * time.Second)
		c.AggrPidList = msg.AggrPidList

		ctx.Request(c.selectorPid, &serverMsg.AssignDevicesToAggr{NumAggr: numAggr,
			DevicesPerAggr: c.devices_for_each_aggr,
			AggrPidList:    msg.AggrPidList})

	case *serverMsg.DevicesAssignmentSuccess:
		log.Printf("COORDINATOR ::::: Received message from Selector ---> Devices Assginment Successful\n")
		time.Sleep(1 * time.Second)
		log.Printf("COORDINATOR ::::: Time to send model to all the Aggregators for training\n")

		for i := 0; i < len(c.AggrPidList); i++ {
			ctx.Request(c.AggrPidList[i], &serverMsg.ModelCarryingMessage{Model: "model.sh"})
			time.Sleep(2 * time.Second)
		}

	case *serverMsg.ModelTrainingCompletion:
		c.AggrPidList = nil
		time.Sleep(1 * time.Second)
		ctx.Stop(ctx.Children()[0])
		abort_status = 0
		log.Printf("###################FL ROUND COMPLETED##############")

	default:
	}
}

func newCoordinatorActor() actor.Actor {
	return &coordinatorActor{devices_needed: 20, devices_for_each_aggr: 5}
}
