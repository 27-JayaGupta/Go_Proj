package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type pong struct{}

type ping struct{}

type pingActor struct {
	pongPid *actor.PID
}

func (p *pingActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case struct{}:
		future := ctx.RequestFuture(p.pongPid, &ping{}, time.Second)
		result, err := future.Result()

		if err != nil {
			log.Print(err.Error())
			return
		}

		log.Printf("Received %#v", result)

		// Never comes here.
		// When the pong actor responds to the sender,
		// the sender is not a ping actor but a future process.
	case *pong:
		log.Printf("Received pong message\n")
	}
}

func main() {
	system := actor.NewActorSystem()

	pongProps := actor.PropsFromFunc(func(ctx actor.Context) {
		switch ctx.Message().(type) {
		case *ping:
			// Below both work in this example, but their behavior slightly differ.
			// ctx.Sender().Tell() panics and recovers if the sender is nil;
			// while ctx.Respond() checks the presence of sender and redirects the message to dead letter process
			// when sender is absent.
			log.Printf("Received ping message\n")
			ctx.Respond(&pong{})

		default:

		}
	})

	pongPid := system.Root.Spawn(pongProps)

	pingProps := actor.PropsFromProducer(func() actor.Actor {
		return &pingActor{
			pongPid: pongPid,
		}
	})

	pingPid := system.Root.Spawn(pingProps)

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, os.Kill)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			system.Root.Send(pingPid, struct{}{})
		case <-finish:
			log.Printf(("Finish\n"))
			return
		}
	}
}
