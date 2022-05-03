package main

import (
	"FL_System/Remote/messages"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
)

type pingActor struct {
	cnt     uint64
	pongPid *actor.PID
}

func (p *pingActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case struct{}:
		p.cnt += 1
		ping := &messages.Ping{
			Cnt: p.cnt,
		}

		ctx.Request(p.pongPid, ping)

	case *messages.Pong:
		log.Printf("Received Pong Message")
	default:
	}
}

func main() {
	system := actor.NewActorSystem()

	config := remote.Configure("127.0.0.1", 8081)
	remoting := remote.NewRemote(system, config)
	remoting.Start()

	remotePong := actor.NewPID("127.0.0.1:8080", "pongActorID")

	pingProps := actor.PropsFromProducer(func() actor.Actor {
		return &pingActor{
			pongPid: remotePong,
		}
	})

	pingPid := system.Root.Spawn(pingProps)

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Kill, os.Interrupt)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			system.Root.Send(pingPid, struct{}{})
		case <-finish:
			log.Printf("Finish\n")
			return
		}
	}
}
