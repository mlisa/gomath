package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/mlisa/gomath/message"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type Coordinator struct {
	MaxPeers     int
	Peers        map[string]*actor.PID
	Coordinators map[string]*actor.PID
	Controller   *Controller
}

var pings map[string]int64
var mutex = &sync.Mutex{}

func (coordinator *Coordinator) Receive(context actor.Context) {
	log := coordinator.Controller.Log
	switch msg := context.Message().(type) {
	case *message.Hello:
		// Coordinator received an 'Hello' message from a want which want to join the region, check availability
		if len(coordinator.Peers) < coordinator.MaxPeers {
			context.Respond(&message.Available{})
		}
		context.Self().Tell(&message.Ping{})
	case *message.Register:
		// Peer wants to be registered in the region, update the nodes
		log(fmt.Sprintf("Added peer '%s' to region", context.Sender().Id))
		context.Sender().Request(&message.Welcome{coordinator.Peers}, context.Self())
		// update all others peers newnode
		for _, PID := range coordinator.Peers {
			actor.NewPID(PID.Address, PID.Id).Tell(&message.NewNode{context.Sender()})
		}
		context.Watch(context.Sender())
		coordinator.Peers[context.Sender().String()] = context.Sender()
		context.Self().Tell(&message.Ping{})

	case *message.RequestForCache:
		// Received a request from an another coordinator to forward to each peer
		log(fmt.Sprintf("Request for '%s' from '%s'", msg.Operation, context.Sender().Id))
		if response := coordinator.sendToAll(context.Self(), coordinator.Coordinators, &message.RequestForCacheExternal{msg.Operation}); response != nil {
			context.Sender().Request(response.(*message.Response), context.Self())
		}
		context.Self().Tell(&message.Ping{})
	case *message.RequestForCacheExternal:
		// Received a request from a peer to forward to each known coordinator
		log(fmt.Sprintf("Request for '%s' from '%s'", msg.Operation, context.Sender().Id))
		if response := coordinator.sendToAll(context.Self(), coordinator.Peers, &message.RequestForCache{msg.Operation}); response != nil {
			context.Respond(response.(*message.Response))
		}
		context.Self().Tell(&message.Ping{})

	case *message.Pong:
		// Received a Pong from a previous Ping from a peer
		mutex.Lock()
		latency := msg.Pong - pings[context.Sender().String()]
		pings[context.Sender().String()] = latency
		coordinator.Controller.UpdatePings(pings)
		mutex.Unlock()
	case *message.Ping:
		// Pinga all pees
		mutex.Lock()
		pings = make(map[string]int64, len(coordinator.Peers))
		if len(coordinator.Peers) > 0 {
			for _, PID := range coordinator.Peers {
				ping := time.Now().UnixNano() / 1000000
				pings[PID.String()] = ping
				actor.NewPID(PID.Address, PID.Id).Tell(&message.Ping{})
			}
		}
		mutex.Unlock()
	case *message.GetPing:
		if ping, ok := pings[msg.Peer]; ok {
			context.Respond(&message.Pong{Pong: ping})
		}
		context.Self().Tell(&message.Ping{})

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
		for _, PID := range coordinator.Peers {
			actor.NewPID(PID.Address, PID.Id).Request(&message.DeadNode{msg.Who}, context.Self())
		}
	}

}

func (c *Coordinator) sendToAll(from *actor.PID, who map[string]*actor.PID, what interface{}) interface{} {
	// Channel to stop all goroutines
	response := make(chan interface{})
	for _, PID := range who {
		if PID.Address != from.Address {
			go func() {
				res, _ := actor.NewPID(PID.Address, PID.Id).RequestFuture(what, 5*time.Second).Result()
				response <- res
			}()
		}
	}
	for range who {
		val := <-response
		if val, ok := val.(*message.Response); ok {
			return val
		}
	}
	return nil
}
