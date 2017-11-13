package main

import (
	"log"
	"time"

	"github.com/mlisa/gomath/message"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type Coordinator struct {
	MaxPeers     int
	Peers        map[string]*actor.PID
	Coordinators map[string]*actor.PID
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
		log.Printf("[COORDINATOR] Added peer '%s' to region", context.Sender().Id)
		context.Sender().Request(&message.Welcome{coordinator.Peers}, context.Self())
		// update all others peers newnode
		for _, PID := range coordinator.Peers {
			actor.NewPID(PID.Address, PID.Id).Tell(&message.NewNode{context.Sender()})
		}
		context.Watch(context.Sender())
		coordinator.Peers[context.Sender().String()] = context.Sender()
	case *actor.Stopping:
		log.Println("[COORDINATOR] Stopping, actor is about shut down")
	case *actor.Stopped:
		log.Println("[COORDINATOR] Stopped, actor and it's children are stopped")
	case *actor.Terminated:
		// Watch for terminated peers of the region
		log.Printf("[COORDINATOR] detected node failure: '%s'", msg.Who.Id)
		if _, present := coordinator.Peers[msg.Who.String()]; present {
			delete(coordinator.Peers, msg.Who.String())
		}
		for _, PID := range coordinator.Peers {
			actor.NewPID(PID.Address, PID.Id).Request(&message.DeadNode{msg.Who}, context.Self())
		}
	case *message.RequestForCache:
		// Received a request from an another coordinator to forward to each peer
		log.Printf("[COORDINATOR] Request for '%s' from '%s'", msg.Operation, context.Sender().Id)
		response := coordinator.sendToAll(context.Self(), coordinator.Coordinators, &message.RequestForCacheExternal{msg.Operation})
		context.Sender().Request(response.(*message.Response), context.Self())
	case *message.RequestForCacheExternal:
		// Received a request from a peer to forward to each known coordinator
		log.Printf("[COORDINATOR] Request for '%s' from '%s'", msg.Operation, context.Sender().Id)
		response := coordinator.sendToAll(context.Self(), coordinator.Peers, &message.RequestForCache{msg.Operation})
		context.Respond(response.(*message.Response))
	}
}

func (c *Coordinator) sendToAll(from *actor.PID, who map[string]*actor.PID, what interface{}) interface{} {
	// Channel to stop all goroutines
	shutdown := make(chan struct{})
	response := make(chan interface{})
	for _, PID := range who {
		if PID.Address != from.Address {
			go func() {
				select {
				default:
					req := actor.NewPID(PID.Address, PID.Id).RequestFuture(what, 5*time.Second)
					if r, err := req.Result(); err == nil {
						response <- r
						close(shutdown)
					}
				case <-shutdown:
					return
				}
			}()
		}
	}
	val := <-response
	return val
}
