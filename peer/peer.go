package main

import (
	"fmt"
	"log"

	"github.com/mlisa/gomath/common"
	"github.com/mlisa/gomath/message"

	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type Peer struct {
	otherNodes  map[string]*actor.PID
	coordinator *actor.PID
	Controller  *Controller
}

func (peer *Peer) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		coordinators, err := common.GetCoordinatorsList() //lettura da file config
		if err == nil {
			log.Println(coordinators)
			for _, PID := range coordinators {
				log.Println("[PEER] Try to connect to " + PID.Address + " " + PID.Id)
				tempCoordinator := actor.NewPID(PID.Address, PID.Id)
				tempCoordinator.Request(&message.Hello{}, context.Self())
			}
		}

	case *message.LostConnectionCoordinator:
		coordinators, err := common.GetCoordinatorsList() //lettura da file config
		if err == nil {
			for _, PID := range coordinators {
				tempCoordinator := actor.NewPID(PID.Address, PID.Id)
				tempCoordinator.Request(&message.Hello{peer.Controller.Config.Myself.Latency, peer.Controller.Config.Myself.ComputationCapability, peer.Controller.Config.Myself.Queue}, context.Self())
			}
		}

	case *message.Available:
		peer.coordinator = context.Sender()
		context.Watch(peer.coordinator)
		peer.Controller.Log(FOUNDNEWCOORDINATOR)
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
	case *message.Ping:
		peer.coordinator.Request(&message.Pong{time.Now().UnixNano() / 1000000}, context.Self())

	case *message.Welcome:
		log.Println("[PEER] I'm in!")
		peer.otherNodes = msg.Nodes
		delete(peer.otherNodes, context.Self().String())
		log.Println(peer.otherNodes)
		context.SetBehavior(peer.Operative)

	case *actor.Stopping:
		fmt.Println("[PEER] Stopping, actor is about shut down")

	case *actor.Stopped:
		fmt.Println("[PEER] Stopped, actor and it's children are stopped")
	}
}

func (peer *Peer) Operative(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.Ping:
		peer.coordinator.Request(&message.Pong{time.Now().UnixNano() / 1000000}, context.Self())

	case *actor.Terminated:
		peer.Controller.Log(LOSTCONNECTION)
		context.SetBehavior(peer.Receive)
		context.Self().Tell(&message.LostConnectionCoordinator{msg.Who})

	case *message.AskForResult:
		res := peer.sendToAll(&message.RequestForCache{Operation: msg.Operation})
		if res == nil || len(peer.otherNodes) == 0 {
			peer.Controller.setLog("No one has the response, contacting coordinator")
			peer.coordinator.Request(&message.RequestForCache{Operation: msg.Operation}, context.Self())
		} else {
			peer.Controller.SetOutput(res.(*message.Response).Result)
		}

	case *message.NewNode:
		peer.Controller.Log(NEWNODE)
		peer.otherNodes[msg.Newnode.String()] = msg.Newnode

	case *message.DeadNode:
		peer.Controller.Log(DEADNODE)
		delete(peer.otherNodes, msg.DeadNode.String())

	case *message.RequestForCache:
		peer.Controller.Log(SEARCHINCACHE)
		res := peer.Controller.SearchInCache(msg.Operation)
		if res != "" {
			peer.Controller.Log(FOUNDRESULTINCACHE)
			context.Respond(&message.Response{Result: res})
		} else {
			peer.Controller.Log(NOTFOUND)
		}

	case *message.Response:
		peer.Controller.SetOutput(msg.Result)

	case *actor.Stopping:
		log.Println("[PEER] Stopping, actor is about shut down")

	case *actor.Stopped:
		log.Println("[PEER] Stopped, actor and it's children are stopped")
	}
}

func (p *Peer) sendToAll(what interface{}) interface{} {
	// Channel to stop all goroutines
	response := make(chan interface{})
	for _, PID := range p.otherNodes {
		go func() {
			p.Controller.setLog("Asking to.." + PID.String())
			req := actor.NewPID(PID.Address, PID.Id).RequestFuture(what, 2*time.Second)
			res, _ := req.Result()
			response <- res
		}()

	}

	for i := 0; i < len(p.otherNodes); i++ {
		val := <-response
		if response, ok := val.(*message.Response); ok {
			return response
		}
	}
	return nil
}
