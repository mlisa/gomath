package peer

import (
	"fmt"
	"log"

	"github.com/mlisa/gomath/common"
	"github.com/mlisa/gomath/message"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type Peer struct {
	otherNodes  []*actor.PID
	coordinator *actor.PID
	Controller  *Controller
}

func (peer *Peer) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		fmt.Println("[PEER] Started, initialize actor here, I'm " + context.Self().Id + " " + context.Self().Address)

		coordinators := common.GetConfig("peer").Coordinators //lettura da file config
		for _, PID := range coordinators {
			log.Println("[PEER] Try to connect to " + PID.Address + " " + PID.Id)
			coord := actor.NewPID(PID.Address, PID.Id)
			coord.Request(&message.Hello{}, context.Self())
		}
	case *message.Available:
		log.Println("[PEER] Found a coordinator!")
		peer.coordinator = context.Sender()
		//TODO: requestfuture se no poi si pianta
		peer.coordinator.Request(&message.Register{}, context.Self())
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
		context.SetBehavior(peer.Operative)
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
		res := peer.Controller.SearchInCache(msg.Operation)
		if res != "" {
			context.Sender().Request(&message.Response{Result: res}, context.Self())
		}
	case *message.Response:
		log.Println("[PEER] Received Response from peer!")
		peer.Controller.SetResult(msg.Result)
	case *actor.Stopping:
		log.Println("[PEER] Stopping, actor is about shut down")
	case *actor.Stopped:
		log.Println("[PEER] Stopped, actor and it's children are stopped")
	}
}
