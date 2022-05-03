package main

import (
	serverMsg "fl_system/messages"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

var abort_status int
var test_mode bool = false

func Server(n int) {

	//abort_status := 0
	system := actor.NewActorSystem()

	var timeBetweenRounds time.Duration = 1000
	coordinatorProps := actor.PropsFromProducer(newCoordinatorActor)
	coordinatorPid := system.Root.Spawn(coordinatorProps)

	selectorProps := actor.PropsFromProducer(newSelectorActor)
	selectorPid := system.Root.Spawn(selectorProps)

	var trainerPid []*actor.PID
	j := 0
	for {
		t := abort_status
		for i := 0; i < n; i++ {
			trainerProps := actor.PropsFromProducer(newTrainer)
			pid, err := system.Root.SpawnNamed(trainerProps, "trainer"+strconv.Itoa(j))
			j++
			if err != nil {
				log.Print(err)
				return
			}
			trainerPid = append(trainerPid, pid)
			time.Sleep(1 * time.Second)
			system.Root.Send(pid, &serverMsg.ServerPingtoTrainer{SelectorPid: selectorPid})
			time.Sleep(1 * time.Second)
		}
		system.Root.Send(selectorPid, &serverMsg.SelectorSendCountToCoordinator{CoordinatorPid: coordinatorPid})
		time.Sleep(5 * time.Second)

		if t < abort_status {
			//Means not enough devices were there....make more devices to ping
			continue
		} else {

			//FL Cycle started it will wait for it to complete
			if test_mode {
				return
			}
			time.Sleep(timeBetweenRounds * time.Second)
		}

	}
}

func main() {
	if test_mode == true {
		var n int
		_, err := fmt.Scanf("%d", &n)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		Server(n)
	} else {
		Server(20)
	}

}
