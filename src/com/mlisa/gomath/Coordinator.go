package main

import (
	"runtime"
	"github.com/AsynkronIT/protoactor-go/remote"
	"com/mlisa/gomath/message"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/goconsole"
	"log"
	"com/mlisa/gomath/Common"
)

var nodes = make(map[string]string)
var coordinatorData = Common.Node{"Coordinator","127.0.0.1:8081"}

func waitingForNodes(context actor.Context){
	switch msg := context.Message().(type) {
	case *message.Hello:
		log.Println("[COORDINATOR] Received from " + msg.Name + " " + msg.Address)
		sender := actor.NewPID(msg.Address, msg.Name)
		sender.Tell(&message.Available{coordinatorData.Address, coordinatorData.Name})

	case *message.Register:
		log.Println("[COORDINATOR] Sending nodes to " + msg.Name + " " + msg.Address)

		sender := actor.NewPID(msg.Address, msg.Name)
		sender.Tell(&message.Welcome{nodes})
		message := &message.NewNode{msg.Address, msg.Name}

		for k, v := range nodes {
			sender = actor.NewPID(v, k)
			sender.Tell(message)
		}
		nodes[msg.Name] = msg.Address

		}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	remote.Start(coordinatorData.Address)

	//remote.Register("Coordinator", actor.FromFunc(waitingForNodes))

	actor.SpawnNamed(actor.FromFunc(waitingForNodes), coordinatorData.Name,)

	console.ReadLine()

}