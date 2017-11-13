package main

import (
	"log"
	"runtime"

	console "github.com/AsynkronIT/goconsole"
	"github.com/mlisa/gomath/common"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	maxpeer = kingpin.Flag("max-peers", "Maximum number of nodes to accept").Short('m').Default("50").Int()
	config  = kingpin.Flag("config", "Configuration file for coordinator").Short('c').Default("config_coordinator.json").String()
	token   = kingpin.Flag("token", "one-time-token to register the coordinator").Short('t').String()
)

func main() {
	kingpin.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	configCoordinator, err := common.GetFileConfig(*config)
	if err != nil {
		kingpin.FatalUsage("Wrong usage, please see the help")
	}

	controller := Controller{}
	if len(*token) > 0 {
		controller.PublishCoordinator(*token)
	}

	if err := controller.StartCoordinator(configCoordinator); err == nil {
		StartGui()
		console.ReadLine()
	} else {
		log.Panicln(err)
	}
}
