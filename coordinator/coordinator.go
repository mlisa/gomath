package main

import (
	"fmt"
	"time"

	"github.com/mlisa/gomath/common"
	"github.com/mlisa/gomath/message"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type Coordinator struct {
	MaxPeers     int
	Peers        map[string]*actor.PID
	Coordinators map[string]*actor.PID
	Controller   *Controller
}

var pings map[string]common.Pong

func (coordinator *Coordinator) Receive(context actor.Context) {
	log := coordinator.Controller.Log
	switch msg := context.Message().(type) {
	case *message.Hello:
		// Coordinator received an 'Hello' message from a want which want to join the region, check availability
		if len(coordinator.Peers) < coordinator.MaxPeers {
			context.Respond(&message.Available{context.Self()})
		}
	case *message.Register:
		// Peer wants to be registered in the region, update the nodes
		if len(coordinator.Peers) < coordinator.MaxPeers {
			log(fmt.Sprintf("Added peer '%s' to region", context.Sender().Id))
			context.Sender().Request(&message.Welcome{coordinator.Peers}, context.Self())
			// update all others peers newnode
			for _, PID := range coordinator.Peers {
				actor.NewPID(PID.Address, PID.Id).Tell(&message.NewNode{context.Sender()})
			}
			context.Watch(context.Sender())
			coordinator.Peers[context.Sender().String()] = context.Sender()
		} else {
			context.Sender().Tell(&message.NotAvailable{})
		}
		context.Self().Tell(&message.Ping{})
		coordinator.Controller.UpdatePings(pings)

	case *message.RequestForCache:
		// Received a request from a peer to forward to each known coordinator
		log(fmt.Sprintf("Request for '%s' from '%s'", msg.Operation, msg.Sender.Id))
		if response := coordinator.sendToAll(context.Self(), coordinator.Coordinators, &message.RequestForCacheExternal{msg.Operation, context.Self()}); response != nil {
			context.Respond(response.(*message.Response))
		}
		context.Self().Tell(&message.Ping{})
	case *message.RequestForCacheExternal:
		// Received a request from an another coordinator to forward to each peer
		log(fmt.Sprintf("Request for '%s' from '%s'", msg.Operation, msg.Sender.Id))
		if response := coordinator.sendToAll(context.Self(), coordinator.Peers, &message.RequestForCache{msg.Operation, context.Self()}); response != nil {
			context.Respond(response.(*message.Response))
		}
		context.Self().Tell(&message.Ping{})
		coordinator.Controller.UpdatePings(pings)

	case *message.Pong:
		// Received a Pong from a previous Ping from a peer
		ping := pings[context.Sender().String()]
		if !ping.Complete {
			pong := time.Now().UnixNano() / 1000000
			pings[context.Sender().String()] = common.Pong{pong - ping.Value, true}
		}
		coordinator.Controller.UpdatePings(pings)
	case *message.Ping:
		// Pings all peers
		if len(pings) <= 0 {
			pings = make(map[string]common.Pong, len(coordinator.Peers))
		}
		if len(coordinator.Peers) > 0 {
			for _, PID := range coordinator.Peers {
				ping := time.Now().UnixNano() / 1000000
				pings[PID.String()] = common.Pong{ping, false}
				actor.NewPID(PID.Address, PID.Id).Tell(&message.Ping{})
			}
		}
	case *message.GetPing:
		if ping, ok := pings[msg.Peer]; ok && ping.Complete {
			context.Respond(&message.Pong{ping.Value})
		}

	case *actor.Stopping:
		log("Stopping, actor is about shut down")
	case *actor.Stopped:
		log("Stopped, actor and it's children are stopped")
	case *actor.Terminated:
		// Watch for terminated peers of the region
		log(fmt.Sprintf("Detected node failure: '%s'", msg.Who.Id))
		if _, present := coordinator.Peers[msg.Who.String()]; present {
			delete(coordinator.Peers, msg.Who.String())
		}
		context.Self().Tell(&message.Ping{})
		coordinator.Controller.UpdatePings(pings)
		for _, PID := range coordinator.Peers {
			actor.NewPID(PID.Address, PID.Id).Request(&message.DeadNode{msg.Who}, context.Self())
		}
	}
}

func (c *Coordinator) sendToAll(from *actor.PID, who map[string]*actor.PID, what interface{}) interface{} {
	// Channel to stop all goroutines
	response := make(chan interface{})
	for _, PID := range who {
		go func(PID *actor.PID) {
			var res interface{}
			if PID.Address != from.Address {
				res, _ = actor.NewPID(PID.Address, PID.Id).RequestFuture(what, 5*time.Second).Result()
			}
			response <- res
		}(PID)
	}
	for range who {
		val := <-response
		if val, ok := val.(*message.Response); ok {
			return val
		}
	}
	return nil
}
