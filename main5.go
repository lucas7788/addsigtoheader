package main

import (
	"github.com/ontio/ontology-eventbus/actor"
	"fmt"
	"time"
	"runtime"
)

type ping struct{ val int }
type pingActor struct{}

func (state *pingActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println("Started, initialize actor here")
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about shut down")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about restart")
	case *ping:
		val := msg.val
		if val < 10 {
			fmt.Println("context.Sender():", context.Sender())
			fmt.Println("context.Self():", context.Self())
			fmt.Println("")
			context.Sender().Request(&ping{val: val + 1}, context.Self())
		} else {
			end := time.Now().UnixNano()
			fmt.Printf("%s end %d\n", context.Self().Id, end)
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	props := actor.FromProducer(func() actor.Actor { return &pingActor{} })
	actora := actor.Spawn(props)
	actorb := actor.Spawn(props)
	fmt.Println("actora:", actora)
	fmt.Println("actorb:", actorb)
	fmt.Println("")
	fmt.Printf("begin time %d\n", time.Now().UnixNano())
	actora.Request(&ping{val: 1}, actorb)
	time.Sleep(10 * time.Second)
	actora.Stop()
	actorb.Stop()
}