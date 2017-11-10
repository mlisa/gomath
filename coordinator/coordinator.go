package main

import (
	"log"

	"github.com/mlisa/gomath/message"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type Coordinator struct {
	MaxPeers     int
	Peers        map[string]*actor.PID
	Coordinators []Coordinators
}

func (coordinator *Coordinator) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.Hello:
		// Coordinator received an 'Hello' message from a want which want to join the region, check availability
		if len(coordinator.Peers) < coordinator.MaxPeers {
			context.Sender().Request(&message.Available{}, context.Self())
		} else {
			context.Sender().Request(&message.NotAvailable{}, context.Self())
		}
	case *message.Register:
		// Peer wants to be registered in the region, update the nodes
		log.Println("[COORDINATOR] Added peer '%s' to region", context.Sender().Id)
		context.Sender().Request(&message.Welcome{coordinator.Peers}, context.Self())
		// update all others peers newnode
		for _, PID := range coordinator.Peers {
			actor.NewPID(PID.Address, PID.Id).Request(&message.NewNode{context.Sender()}, context.Self())
		}
		context.Watch(context.Sender())
		coordinator.Peers[context.Sender().String()] = context.Sender()
		log.Println(context.Sender().String())
	case *actor.Stopping:
		log.Println("[COORDINATOR] Stopping, actor is about shut down")
	case *actor.Stopped:
		log.Println("[COORDINATOR] Stopped, actor and it's children are stopped")
	case *actor.Terminated:
		// Watch for terminated peers of the region
		log.Println("[COORDINATOR] detected node failure: '%s'", msg.Who.Id)
		if _, present := coordinator.Peers[msg.Who.String()]; present {
			delete(coordinator.Peers, msg.Who.String())
		}
		for _, PID := range coordinator.Peers {
			actor.NewPID(PID.Address, PID.Id).Request(&message.DeadNode{msg.Who}, context.Self())
		}
	}
}
