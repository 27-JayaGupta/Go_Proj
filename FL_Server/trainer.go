package main

import (
	serverMsg "fl_system/messages"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type Trainer struct {
	selectorPid *actor.PID
	model       string
}

func (t *Trainer) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		log.Print("[+]TRAINER ::::: Trainer started\n")

	case *serverMsg.ServerPingtoTrainer:
		log.Printf("[+]TRAINER ::::: Server pinged to send availability status to Selector\n")
		t.selectorPid = msg.SelectorPid
		ctx.Request(t.selectorPid, &serverMsg.TrainerPingtoSelector{})

	case *serverMsg.Demo:
		log.Printf("[+]TRAINER ::::: Assigned to Aggregator %v", ctx.Sender().Id)

	case *serverMsg.ModelCarryingMessage:
		log.Printf("[+]TRAINER ::::: %v Reporting.... Received Model %s from Aggregator %v\n", ctx.Self().Id, msg.Model, ctx.Sender().Id)

		log.Printf("Executing model\n")
		cmd := exec.Command("./" + msg.Model)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

		time.Sleep(2 * time.Second)
		log.Printf("Trainer" + ctx.Self().Id + "reporting, Model Training Completed")
		ctx.Respond(&serverMsg.ModelTrainingCompletion{})
	}

}

func newTrainer() actor.Actor {
	return &Trainer{model: ""}
}
