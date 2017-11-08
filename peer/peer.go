package main

import (
	"fmt"
	"log"

	"github.com/mlisa/gomath/common"
	"github.com/mlisa/gomath/message"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type Peer struct {
	otherNodes  map[string]*actor.PID
	coordinator *actor.PID
	Controller  *Controller
	Config      common.Config
}

func (peer *Peer) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		coordinators := peer.Config.Coordinators //lettura da file config
		for _, PID := range coordinators {
			log.Println("[PEER] Try to connect to " + PID.Address + " " + PID.Id)
			tempCoordinator := actor.NewPID(PID.Address, PID.Id)
			log.Println(tempCoordinator.Id)
			tempCoordinator.Request(&message.Hello{}, context.Self())
		}
	case *message.LostConnectionCoordinator:
		coordinators := peer.Config.Coordinators //lettura da file config
		for i, PID := range coordinators {
			if PID.Address == msg.Coordinator.Address {
				log.Println("[PEER] Removing dead coordinator..")
				coordinators = append(coordinators[:i], coordinators[i+1:]...)
			} else {
				log.Println("[PEER] Try to reconnect to " + PID.Address + " " + PID.Id)
				tempCoordinator := actor.NewPID(PID.Address, PID.Id)
				tempCoordinator.Request(&message.Hello{peer.Config.Myself.Latency, peer.Config.Myself.ComputationCapability, peer.Config.Myself.Queue}, context.Self())
			}
		}
	case *message.Available:
		peer.coordinator = context.Sender()
		context.Watch(peer.coordinator)
		log.Println("[PEER] Found a coordinator! " + peer.coordinator.Address + peer.coordinator.Id)
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
	case *actor.Terminated:
		log.Println("[PEER] Lost Connection from coordinator " + msg.Who.Address + msg.Who.Id)
		context.SetBehavior(peer.Receive)
		context.Self().Tell(&message.LostConnectionCoordinator{msg.Who})
	case *message.AskForResult:
		log.Println("[PEER] Sending RequestForCache")
		log.Println(peer.otherNodes)
		for _, peer := range peer.otherNodes {
			log.Println("[PEER] Sending RequestForCache to" + peer.Id + peer.Address)
			peer.Request(&message.RequestForCache{Operation: msg.Operation}, context.Self())
		}
	case *message.NewNode:
		peer.otherNodes[msg.Newnode.Address+msg.Newnode.Id] = msg.Newnode
	case *message.DeadNode:
		delete(peer.otherNodes, msg.DeadNode.Address+msg.DeadNode.Id)
		//peer.otherNodes = append(peer.otherNodes, msg.DeadNode)
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
