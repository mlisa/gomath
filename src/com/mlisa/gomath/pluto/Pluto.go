package pluto

import (
	"runtime"
	"github.com/AsynkronIT/protoactor-go/remote"
	"com/mlisa/gomath/message"
	"github.com/AsynkronIT/protoactor-go/actor"
	"log"
	"github.com/AsynkronIT/goconsole"
)

var otherNodes map[string]string

func reactive (context actor.Context) {
	log.Println("Become reactive")
	switch msg := context.Message().(type) {
	case *message.Welcome:
		log.Println("Connected")
		otherNodes := msg.Nodes
		context.SetBehavior(operational)
	}
}

func operational(context actor.Context){
	log.Println("Become operational")
	switch msg := context.Message().(type) {
	case *message.NewNode:
		println("New node arrived")
		otherNodes[msg.Name] = msg.Address
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	remote.Start("127.0.0.1:8083")
	//remote.Register("Pippo", actor.FromFunc(reactive))

	actor.SpawnNamed( actor.FromFunc(reactive), "Pippa")

	message := &message.Hello{
		Name: "Pippa",
		Address: "127.0.0.1:8083",
	}

	coordinator := actor.NewPID("127.0.0.1:8081", "Coordinator")

	log.Println("Sending hello")
	coordinator.Tell(message)

	console.ReadLine()

}