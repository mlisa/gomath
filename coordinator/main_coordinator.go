package main

import (
	"log"
	"runtime"

	"github.com/mlisa/gomath/coordinator"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	id := "coordinator"
	address := "127.0.0.1:8081"
	remote.Start(address)

	props := actor.FromInstance(&coordinator.Coordinator{MaxPeers: 50, Peers: make([]*actor.PID, 0, 50)})
	_, err := actor.SpawnNamed(props, id)
	if err != nil {
		log.Panicln(err)
	}
	console.ReadLine()
}
