package coordinator

import (
	"com/mlisa/gomath/message"
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type Coordinator struct {
	MaxPeers int
	Peers    []*actor.PID
}

func (coordinator *Coordinator) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
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
		// update all others peers newnode
		for _, PID := range coordinator.Peers {
			actor.NewPID(PID.Address, PID.Id).Request(&message.NewNode{context.Sender()}, context.Self())
		}
		context.Watch(context.Sender())
		coordinator.Peers = append(coordinator.Peers, context.Sender())
	case *actor.Stopping:
		log.Println("[COORDINATOR] Stopping, actor is about shut down")
	case *actor.Stopped:
		log.Println("[COORDINATOR] Stopped, actor and it's children are stopped")
	case *actor.Terminated:
		log.Println(msg)
	}
}
