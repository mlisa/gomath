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
			tempCoordinator.Request(&message.Hello{}, context.Self())
		}

	case *message.LostConnectionCoordinator:
		for i, PID := range peer.Config.Coordinators {
			if PID.Address == msg.Coordinator.Address {
				log.Println("[PEER] Removing dead coordinator..")
				peer.Config.Coordinators = append(peer.Config.Coordinators[:i], peer.Config.Coordinators[i+1:]...)
			} else {
				tempCoordinator := actor.NewPID(PID.Address, PID.Id)
				tempCoordinator.Request(&message.Hello{peer.Config.Myself.Latency, peer.Config.Myself.ComputationCapability, peer.Config.Myself.Queue}, context.Self())
			}
		}

	case *message.Available:
		peer.coordinator = context.Sender()
		context.Watch(peer.coordinator)
		peer.Controller.Log(FOUNDNEWCOORDINATOR)
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

	case *actor.Terminated:
		peer.Controller.Log(LOSTCONNECTION)
		context.SetBehavior(peer.Receive)
		context.Self().Tell(&message.LostConnectionCoordinator{msg.Who})

	case *message.AskForResult:
		for _, otherPeer := range peer.otherNodes {
			if otherPeer.Id != context.Self().Id || otherPeer.Address != context.Self().Address {
				peer.Controller.setLog("Asking to.." + otherPeer.String())
				otherPeer.Request(&message.RequestForCache{Operation: msg.Operation}, context.Self())
			}
		}
		context.SetBehavior(peer.WaitingForResponse)
	case *message.NewNode:
		peer.Controller.Log(NEWNODE)
		peer.otherNodes[msg.Newnode.String()] = msg.Newnode

	case *message.DeadNode:
		peer.Controller.Log(DEADNODE)
		delete(peer.otherNodes, msg.DeadNode.String())

	case *message.RequestForCache:
		peer.Controller.setLog("Received RequestForCache message from peer" + context.Sender().String())
		res := peer.Controller.SearchInCache(msg.Operation)
		if res != "" {
			context.Sender().Tell(&message.Response{Result: res})
		} else {
			context.Sender().Tell(&message.NotFound{msg.Operation})
		}

	case *actor.Stopping:
		log.Println("[PEER] Stopping, actor is about shut down")

	case *actor.Stopped:
		log.Println("[PEER] Stopped, actor and it's children are stopped")
	}
}

func (peer *Peer) WaitingForResponse(context actor.Context) {
	numResponse := 0
	switch msg := context.Message().(type) {
	case *message.Response:
		peer.Controller.Log(RECEIVEDRESPONSE)
		peer.Controller.SetOutput(msg.Result)
		context.SetBehavior(peer.Operative)
	case *message.NotFound:
		numResponse++
		if numResponse == len(peer.otherNodes) {
			peer.coordinator.Request(&message.RequestForCache{Operation: msg.Operation}, context.Self())
		} else if numResponse == len(peer.otherNodes)+1 {
			//Forza il calcolo in locale
			context.SetBehavior(peer.Operative)
		}
	}
}
