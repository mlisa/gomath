package main

import (
	"log"
	"runtime"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/mlisa/gomath/common"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	maxpeer = kingpin.Flag("max-peers", "Maximum number of nodes to accept").Short('m').Default("50").Int()
	config  = kingpin.Flag("config", "Configuration file for coordinator").Short('c').Default("config_coordinator.json").String()
)

func main() {
	kingpin.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	configCoordinator, err := common.GetFileConfig(*config)
	if err != nil {
		kingpin.FatalUsage("Wrong usage, please see the help")
	}
	remote.Start(configCoordinator.Myself.Address)

	props := actor.FromInstance(&Coordinator{MaxPeers: *maxpeer, Peers: make([]*actor.PID, 0, *maxpeer)})
	_, err = actor.SpawnNamed(props, configCoordinator.Myself.Id)
	if err != nil {
		log.Panicln(err)
	}
	console.ReadLine()
}
