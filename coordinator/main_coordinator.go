package main

import (
	"log"
	"runtime"

	"github.com/mlisa/gomath/common"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	maxpeer = kingpin.Flag("max-peers", "Maximum number of nodes to accept").Short('m').Default("2").Int()
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

	gui := &GuiCoordinator{}
	controller := &Controller{Gui: gui}
	if len(*token) > 0 {
		controller.PublishCoordinator(*token)
	}

	if err := controller.StartCoordinator(configCoordinator); err == nil {
		gui.StartGui(controller)
	} else {
		log.Panicln(err)
	}
}
