package main
/*
import (
	"runtime"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/AsynkronIT/protoactor-go/actor"

)


func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	remote.Start("127.0.0.1:8080")

	//create an actor receiving messages and pushing them onto the channel
	props := actor.FromFunc(Receive)

	pid, err := actor.SpawnNamed(props, "Peer")

	if(err != nil){
		println("Name already in use")
	}

	controller := Controller{Peer: pid}

	controller.Start()

}
*/