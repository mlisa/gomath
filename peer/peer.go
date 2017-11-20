package main

import (
	"fmt"
	"log"

	"github.com/mlisa/gomath/common"
	"github.com/mlisa/gomath/message"

	"time"

	"math/rand"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type Peer struct {
	otherNodes        map[string]*actor.PID
	coordinator       *actor.PID
	Controller        *Controller
	computeCapability float32
}

func (peer *Peer) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		response := peer.lookForCoordinator(nil)
		if response != nil {
			peer.coordinator = response.Sender
			peer.coordinator.Request(&message.Register{peer.computeCapability}, context.Self())
			context.SetBehavior(peer.Connected)
		}
	case *message.LookForCoordinator:
		response := peer.lookForCoordinator(nil)

		peer.coordinator = response.Sender
		peer.coordinator.Request(&message.Register{peer.computeCapability}, context.Self())
		context.SetBehavior(peer.Connected)

	case *message.LostConnectionCoordinator:

		response := peer.lookForCoordinator(msg.Coordinator)

		peer.coordinator = response.Sender
		peer.coordinator.Request(&message.Register{peer.computeCapability}, context.Self())
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
		r := rand.Intn(500)
		time.Sleep(time.Millisecond * time.Duration(r))
		peer.coordinator.Request(&message.Pong{time.Now().UnixNano() / 1000000}, context.Self())

	case *message.Welcome:
		peer.Controller.Log(FOUNDNEWCOORDINATOR, peer.coordinator.String())
		context.Watch(peer.coordinator)
		peer.otherNodes = msg.Nodes
		delete(peer.otherNodes, context.Self().String())
		context.SetBehavior(peer.Operative)

	case *message.NotAvailable:
		context.SetBehavior(peer.Receive)
		context.Self().Tell(&message.LookForCoordinator{})
	case *actor.Stopping:
		fmt.Println("[PEER] Stopping, actor is about shut down")

	case *actor.Stopped:
		fmt.Println("[PEER] Stopped, actor and it's children are stopped")
	}
}

func (peer *Peer) Operative(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.Ping:
		r := rand.Intn(500)
		time.Sleep(time.Millisecond * time.Duration(r))
		peer.coordinator.Request(&message.Pong{time.Now().UnixNano() / 1000000}, context.Self())

	case *actor.Terminated:
		peer.Controller.Log(LOSTCONNECTION, msg.Who.String())
		context.SetBehavior(peer.Receive)
		context.Self().Tell(&message.LostConnectionCoordinator{msg.Who})

	case *message.AskForResult:
		res := common.SendToAll(context.Self(), peer.otherNodes, &message.RequestForCache{msg.Operation, context.Self()})
		if res == nil || len(peer.otherNodes) == 0 {
			peer.Controller.Log(ASKCOORDINATOR, peer.coordinator.String())
			future := peer.coordinator.RequestFuture(&message.RequestForCache{msg.Operation, context.Self()}, 6*time.Second)
			r, err := future.Result()
			if err != nil {
				peer.Controller.Log(NORESPONSE, "")
				peer.Controller.ComputeLocal(msg.Operation)
			} else {
				peer.Controller.Log(EXTERNALANSWER, "")
				peer.Controller.SetOutput(r.(*message.Response).Result)
			}
		} else {
			peer.Controller.SetOutput(res.(*message.Response).Result)
		}

	case *message.NewNode:
		peer.Controller.Log(NEWNODE, msg.Newnode.String())
		peer.otherNodes[msg.Newnode.String()] = msg.Newnode

	case *message.DeadNode:
		peer.Controller.Log(DEADNODE, msg.DeadNode.String())
		delete(peer.otherNodes, msg.DeadNode.String())

	case *message.RequestForCache:
		peer.Controller.Log(SEARCHINCACHE, msg.Sender.String())
		res := peer.Controller.SearchInCache(msg.Operation)
		if res != "" {
			peer.Controller.Log(FOUNDRESULTINCACHE, "")
			context.Respond(&message.Response{Result: res})
		} else {
			peer.Controller.Log(NOTFOUND, "")
		}

	case *message.Response:
		peer.Controller.SetOutput(msg.Result)

	case *actor.Stopping:
		log.Println("[PEER] Stopping, actor is about shut down")

	case *actor.Stopped:
		log.Println("[PEER] Stopped, actor and it's children are stopped")
	}
}

func (peer *Peer) lookForCoordinator(deadCoordinator *actor.PID) *message.Available {
	coordinators, err := common.GetCoordinatorsList() //lettura da file config
	if err == nil {
		coordChannel := make(chan interface{})
		for _, PID := range coordinators {
			if deadCoordinator == nil || (deadCoordinator != nil && !PID.Equal(deadCoordinator)) {
				go func(PID *actor.PID) {
					tempCoordinator := actor.NewPID(PID.Address, PID.Id)
					fut := tempCoordinator.RequestFuture(&message.Hello{}, 3*time.Second)
					res, err := fut.Result()
					if err == nil {
						coordChannel <- res
					}
				}(PID)
			}
		}

		val := <-coordChannel
		if response, ok := val.(*message.Available); ok {
			return response
		}
	}
	peer.Controller.Connected = false
	return nil
}
