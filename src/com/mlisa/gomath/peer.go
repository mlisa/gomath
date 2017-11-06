package main

import (
	"com/mlisa/gomath/common"
	"com/mlisa/gomath/message"
	"fmt"
	"log"
	"runtime"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

type Peer struct {
	otherNodes  []*actor.PID
	coordinator *actor.PID
	controller  *actor.PID
}

func (peer *Peer) Receive(context actor.Context) {

	switch context.Message().(type) {
	case *actor.Started:
		fmt.Println("[PEER] Started, initialize actor here, I'm " + context.Self().Id + " " + context.Self().Address)

		coordinators := common.GetConfig().Coordinators //lettura da file config
		for _, PID := range coordinators {
			log.Println("[PEER] Try to connect to " + PID.Address + " " + PID.Id)
			coordinator := actor.NewPID(PID.Address, PID.Id)
			coordinator.Request(&message.Hello{}, context.Self())
		}

	case *message.Available:
		log.Println("[PEER] Found a coordinator!")
		peer.coordinator = context.Sender()
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
		otherNodes := msg.Nodes
		log.Println(otherNodes)
		context.SetBehavior(peer.Operative)
	case *actor.Stopping:
		fmt.Println("[PEER] Stopping, actor is about shut down")
	case *actor.Stopped:
		fmt.Println("[PEER] Stopped, actor and it's children are stopped")
	}
}

func (peer *Peer) Operative(context actor.Context) {
	switch context.Message().(type) {
	case *message.RequestForCache:
		log.Println("[PEER] request for cached result")
	case *actor.Stopping:
		fmt.Println("[PEER] Stopping, actor is about shut down")
	case *actor.Stopped:
		fmt.Println("[PEER] Stopped, actor and it's children are stopped")
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	remote.Start(common.GetConfig("peer").Myself.Address)
	//create an actor receiving messages and pushing them onto the channel
	props := actor.FromInstance(&Peer{})
	_, err := actor.SpawnNamed(props, common.GetConfig("peer").Myself.Id)

	if err != nil {
		println("[PEER] Name already in use")
	}
	console.ReadLine()
}
