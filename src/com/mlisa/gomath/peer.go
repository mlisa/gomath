package main

import (
	"com/mlisa/gomath/common"
	"com/mlisa/gomath/message"
	"fmt"
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type Peer struct {
	otherNodes []*actor.PID

	coordinator *actor.PID

	controller *actor.PID

	receivedNodes bool

	connectedToController bool
}

func (peer *Peer) Receive(context actor.Context) {

	switch context.Message().(type) {
	case *actor.Started:
		fmt.Println("[PEER] Started, initialize actor here, I'm " + context.Self().Id + " " + context.Self().Address)

		coordinators := common.GetConfig().Coordinators //lettura da file config
		for _, PID := range coordinators {
			log.Println("[PEER] Try to connect to " + PID.Address + " " + PID.Name)
			coord := actor.NewPID(PID.Address, PID.Name)
			coord.Request(&message.Hello{context.Self()}, context.Self())
		}
	case *message.Hello:
		peer.connectedToController = true
		peer.controller = context.Sender()
		log.Println(peer.controller.Id)
	case *message.Available:
		log.Println("[PEER] Found a coordinator!")
		peer.coordinator = context.Sender()
		peer.coordinator.Request(&message.Register{context.Self()}, context.Self())
		context.SetBehavior(peer.Connected)

	case *actor.Stopping:
		fmt.Println("[PEER] Stopping, actor is about shut down")
	case *actor.Stopped:
		fmt.Println("[PEER] Stopped, actor and it's children are stopped")
	}

}

func (peer *Peer) Connected(context actor.Context) {
	switch msg := context.Message().(type) {

	case *message.Welcome:
		log.Println("[PEER] I'm in!")
		peer.otherNodes = msg.Nodes
		peer.receivedNodes = true
		if peer.connectedToController && peer.receivedNodes {
			context.SetBehavior(peer.Operative)
		}
	case *message.Hello:
		peer.controller = context.Sender()
		log.Println(peer.controller.Id)
		peer.connectedToController = true
		if peer.connectedToController && peer.receivedNodes {
			context.SetBehavior(peer.Operative)
		}
	case *actor.Stopping:
		fmt.Println("[PEER] Stopping, actor is about shut down")
	case *actor.Stopped:
		fmt.Println("[PEER] Stopped, actor and it's children are stopped")
	}
}

func (peer *Peer) Operative(context actor.Context) {
	switch msg := context.Message().(type) {

	case *message.AskForResult:
		log.Println("[PEER] Sending RequestForCache")
		log.Println(peer.otherNodes)
		for _, peer := range peer.otherNodes {
			log.Println("[PEER] Sending RequestForCache to" + peer.Id + peer.Address)
			peer.Request(&message.RequestForCache{Operation: msg.Operation}, context.Self())
		}
	case *message.RequestForCache:
		log.Println("[PEER] Received RequestForCache")
		log.Println(peer.controller)
		peer.controller.Request(&message.SearchInCache{msg.Operation, context.Sender()}, context.Self())
	case *message.ResponseFromCache:
		log.Println("[PEER] Sending ResponseFromCache")
		msg.SendTo.Request(&message.Response{Result: msg.Result}, context.Self())

	case *message.Response:
		log.Println("[PEER] Received Response from peer!")
		peer.controller.Tell(msg)
	case *actor.Stopping:
		log.Println("[PEER] Stopping, actor is about shut down")
	case *actor.Stopped:
		log.Println("[PEER] Stopped, actor and it's children are stopped")
	}
}
