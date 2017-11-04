package main

import (
	"com/mlisa/gomath/message"
	"log"
	"runtime"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

type coordinatorInfo struct {
	PID     *actor.PID
	Name    string
	Address string
}

// Max 50 nodes per region
var maxNodes = 50
var nodes = make([]*actor.PID, 0, maxNodes)
var coordinator coordinatorInfo

func waitingForNodes(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.Hello:
		log.Println("[COORDINATOR] message \"Hello\" from " + msg.Sender.Id + " " + msg.Sender.Address)
		/// check availability
		context.Sender().Request(&message.Available{coordinator.PID}, context.Self())
	case *message.Register:
		log.Println("[COORDINATOR] Sending region nodes to " + msg.Sender.Id + " " + msg.Sender.Address)

		//sender := actor.NewPID(msg.Address, msg.Name)
		msg.Sender.Tell(&message.Welcome{nodes})
		//message := &message.NewNode{msg.Address, msg.Name}

		for range nodes {
			//	sender = actor.NewPID(v, k)
			//sender.Tell(message)
		}
		//nodes[msg.Name] = msg.Address

	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	name := "coordinator"
	address := "127.0.0.1:8081"
	remote.Start(address)
	props := actor.FromFunc(waitingForNodes)
	pid, err := actor.SpawnNamed(props, name)
	if err != nil {
		log.Panicln(err)
	}
	coordinator = coordinatorInfo{pid, name, address}

	console.ReadLine()
}
