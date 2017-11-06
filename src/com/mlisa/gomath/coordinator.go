package main

import (
	"com/mlisa/gomath/common"
	"com/mlisa/gomath/message"
	"log"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

type Coordinator struct {
	MaxPeers int
	Peers    []*actor.PID
}

func (coordinator *Coordinator) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *message.Hello:
		log.Println("[COORDINATOR] message \"Hello\" from peer " + context.Sender().Id)
		if len(coordinator.Peers) < coordinator.MaxPeers {
			context.Sender().Request(&message.Available{}, context.Self())
		} else {
			context.Sender().Request(&message.NotAvailable{}, context.Self())
		}
	case *message.Register:
		log.Println("[COORDINATOR] Sending {{region}} nodes to " + context.Sender().Id + " " + context.Sender().Address)
		context.Sender().Request(&message.Welcome{coordinator.Peers}, context.Self())
		coordinator.Peers = append(coordinator.Peers, context.Sender())
		// update all others peers newnode
		for _, PID := range coordinator.Peers {
			actor.NewPID(PID.Address, PID.Id).Request(&message.NewNode{context.Sender()}, context.Self())
		}
	case *actor.Stopping:
		log.Println("[COORDINATOR] Stopping, actor is about shut down")
	case *actor.Stopped:
		log.Println("[COORDINATOR] Stopped, actor and it's children are stopped")
	}
}

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	remote.Start(common.GetConfig("coordinator").Myself.Address)
	props := actor.FromInstance(&Coordinator{MaxPeers: 50, Peers: make([]*actor.PID, 0, 50)})
	_, err := actor.SpawnNamed(props, common.GetConfig("coordinator").Myself.Id)
	if err != nil {
		log.Panicln(err)
	}
	console.ReadLine()
}
