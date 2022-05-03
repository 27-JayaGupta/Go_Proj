package main

import (
	serverMsg "fl_system/messages"
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type selectorActor struct {
	trainers          []*actor.PID
	devices_available int
}

func (s *selectorActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *serverMsg.SelectorSendCountToCoordinator:
		log.Printf("############# FL CYCLE STARTED ########\n")
		log.Print("[+] SELECTOR ::::: Invoked by server to send devices available count to Coordinator\n")
		time.Sleep(1 * time.Second)
		ctx.Request(msg.CoordinatorPid, &serverMsg.SelectorInfoMsg{Devices_available: s.devices_available})

	case *serverMsg.NotEnoughDevicesAvailable:
		log.Print("[+]  SELECTOR ::::: Not enough devices available , waiting for more devices to join\n")
		log.Printf("#### TERMINATING FL ROUND .... WAITING FOR MORE DEVICES TO JOIN  #### \n")
		abort_status += 1
		// log.Printf("Abort Status Selector: %d", abort_status)
		return

	case *serverMsg.AssignDevicesToAggr:
		j := 0
		for i := 0; i < msg.NumAggr; i++ {
			log.Printf("[+] SELECTOR ::::: Assigning devices to Aggregator %v", msg.AggrPidList[i])
			var devicesList []*actor.PID
			for k := 0; k < msg.DevicesPerAggr; k++ {
				devicesList = append(devicesList, s.trainers[j])
				j++
			}
			time.Sleep(1 * time.Second)
			ctx.Request(msg.AggrPidList[i], &serverMsg.ReceiveDevices{DevicesPidList: devicesList})
		}

		time.Sleep(2 * time.Second)
		ctx.Respond(&serverMsg.DevicesAssignmentSuccess{})

	case *serverMsg.TrainerPingtoSelector:
		log.Printf("[+] SELECTOR ::::: Connection received from device :%#v", ctx.Sender().Id)
		s.trainers = append(s.trainers, ctx.Sender())
		s.devices_available += 1

	default:
	}
}

func newSelectorActor() actor.Actor {
	return &selectorActor{devices_available: 0}
}
