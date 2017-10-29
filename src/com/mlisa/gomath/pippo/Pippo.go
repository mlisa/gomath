package pippo

import (
	"runtime"
	"github.com/AsynkronIT/protoactor-go/remote"
	"com/mlisa/gomath/message"
	"github.com/AsynkronIT/protoactor-go/actor"
	"log"
	"github.com/AsynkronIT/goconsole"
)


func reactive (context actor.Context) {
	log.Println("Become reactive")
	switch msg := context.Message(); msg.(type) {
	case *message.Ok :
		log.Println("Connected")
		context.SetBehavior(operational)
	}
}

func operational(context actor.Context){
	log.Println("Become operational")
	switch msg := context.Message(); msg.(type) {
	case *message.NewNode:
		println("New node arrived")
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	remote.Start("127.0.0.1:8082")
	//remote.Register("Pippo", actor.FromFunc(reactive))

	actor.SpawnNamed( actor.FromFunc(reactive), "Pippo")

	message := &message.Hello{
		Name: "Pippo",
		Address: "127.0.0.1:8082",
	}

	coordinator := actor.NewPID("127.0.0.1:8081", "Coordinator")

	log.Println("Sending hello")
	coordinator.Tell(message)

	console.ReadLine()

}