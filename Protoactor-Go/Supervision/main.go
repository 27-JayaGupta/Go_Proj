package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type hello struct{ name string }

type childActor struct{}

func (c *childActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *hello:
		log.Printf("Hello %+v", msg.name)
	case *actor.Started:
		log.Printf("Actor has started,Initialize actor here\n")
	case *actor.Stopping:
		log.Printf("Stopping,actor is about to shut down\n")
	case *actor.Restarting:
		log.Printf("Restarting Actor\n")
	case *actor.Stopped:
		log.Printf("Actor and all its children are stopped\n")
	default:

	}
}

func newChildActor() actor.Actor {
	return &childActor{}
}

type rootActor struct {
	childPid_ []*actor.PID
}

func (r *rootActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *hello:
		props := actor.PropsFromProducer(newChildActor)
		log.Print(r.childPid_)
		childPid := ctx.Spawn(props)
		r.childPid_ = append(r.childPid_, childPid)
		ctx.Send(childPid, msg)
	default:

	}
}

func newRootActor() actor.Actor {
	return &rootActor{}
}

func main() {
	system := actor.NewActorSystem()

	decider := func(reason interface{}) actor.Directive {
		log.Printf("Actor failure handling,reason %v\n", reason)
		return actor.StopDirective
	}

	supervisor := actor.NewOneForOneStrategy(10, 1000, decider)

	rootProps := actor.PropsFromProducer(newRootActor, actor.WithSupervisor(supervisor))
	rootPid := system.Root.Spawn(rootProps)

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Kill, os.Interrupt)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			system.Root.Send(rootPid, &hello{name: "Jaya"})

		case <-finish:
			log.Printf("Finish\n")
			return
		}
	}

}
