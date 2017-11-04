package main

import (
	"com/mlisa/gomath/common"
	"com/mlisa/gomath/message"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/AsynkronIT/protoactor-go/actor"
)

var myself *actor.PID

var otherNodes = make([]*actor.PID, 0)

var operationsDone map[string]string

var coord *actor.PID

func getConfig() common.Config {
	absPath, _ := filepath.Abs("config.json")
	file, err := os.Open(absPath)
	if err != nil {
		log.Println("[ERROR] " + err.Error())
	}
	decoder := json.NewDecoder(file)
	configuration := common.Config{}
	decoder.Decode(&configuration)
	return configuration
}

func Receive(context actor.Context) {

	switch msg := context.Message().(type) {
	case *actor.Started:
		myself = actor.NewPID(getConfig().Myself.Address, getConfig().Myself.Name)
		fmt.Println("[PEER] Started, initialize actor here, I'm " + myself.Id + " " + myself.Address)

		coordinators := getConfig().Coordinators //lettura da file config
		for _, PID := range coordinators {
			log.Println("[PEER] Try to connect to " + PID.Address + " " + PID.Name)
			coordinator := actor.NewPID(PID.Address, PID.Name)
			coordinator.Tell(&message.Hello{myself})
		}

	case *message.Available:
		log.Println("[PEER] Found a coordinator! " + myself.Id + " " + myself.Address)
		coord = msg.Sender
		coord.Tell(&message.Register{myself})
		context.SetBehavior(Connected)

	case *actor.Stopping:
		fmt.Println("[PEER] Stopping, actor is about shut down")
	case *actor.Stopped:
		fmt.Println("[PEER] Stopped, actor and it's children are stopped")
	}

}

func Connected(context actor.Context) {
	switch msg := context.Message().(type) {

	case *message.Welcome:
		log.Println("[PEER] I'm in!")
		otherNodes := msg.Nodes
		log.Println(otherNodes)
		context.SetBehavior(Operative)
	case *actor.Stopping:
		fmt.Println("[PEER] Stopping, actor is about shut down")
	case *actor.Stopped:
		fmt.Println("[PEER] Stopped, actor and it's children are stopped")
	}
}

func Operative(context actor.Context) {
	switch msg := context.Message().(type) {

	case *message.RequestForCache:
		result, doesExist := operationsDone[msg.Operation]
		if doesExist {
			response := &message.Response{myself, result}
			msg.Sender.Tell(response)
		}

	case *actor.Stopping:
		fmt.Println("[PEER] Stopping, actor is about shut down")
	case *actor.Stopped:
		fmt.Println("[PEER] Stopped, actor and it's children are stopped")
	}
}

/*
func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	remote.Start(myself.Address)

	//create an actor receiving messages and pushing them onto the channel
	props := actor.FromFunc(Receive)

	myself, err := actor.SpawnNamed(props, getConfig().Myself.Name)
	log.Println(myself.Id)
	if err != nil {
		println("[PEER] Name already in use")
	}

	console.ReadLine()

}
*/
