package main

import (
	"com/mlisa/gomath/Common"
	"com/mlisa/gomath/message"
	"encoding/json"
	"fmt"
	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var myself Common.PID

var otherNodes map[string]string

var operationsDone map[string]string

var coordinator *actor.PID

func getConfig() Common.Config {
	absPath, _ := filepath.Abs("com/mlisa/gomath/config.json")
	file, err := os.Open(absPath)
	if err != nil {
		log.Println("[ERROR] " + err.Error())
	}
	decoder := json.NewDecoder(file)
	configuration := Common.Config{}
	decoder.Decode(&configuration)
	return configuration
}

func Receive(context actor.Context) {

	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println("[PEER] Started, initialize actor here")

		coordinators := getConfig().Coordinators //lettura da file config
		for _, PID := range coordinators {
			log.Println("[PEER] Try to connect to " + PID.Address + " " + PID.Name)
			coordinator := actor.NewPID(PID.Address, PID.Name)
			coordinator.Tell(&message.Hello{myself.Address, myself.Name})
		}

	case *message.Available:
		log.Println("[PEER] Found a coordinator!")
		coordinator = actor.NewPID(msg.Address, msg.Name)
		coordinator.Tell(&message.Register{myself.Address, myself.Name})
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
			response := &message.Response{result, myself.Name, myself.Name}
			sender := actor.NewPID(msg.SenderAddress, msg.SenderName)
			sender.Tell(response)
		}

	case *actor.Stopping:
		fmt.Println("[PEER] Stopping, actor is about shut down")
	case *actor.Stopped:
		fmt.Println("[PEER] Stopped, actor and it's children are stopped")
	}
}

func main() {
	myself = getConfig().Myself
	runtime.GOMAXPROCS(runtime.NumCPU())
	remote.Start(myself.Address)

	//create an actor receiving messages and pushing them onto the channel
	props := actor.FromFunc(Receive)

	_, err := actor.SpawnNamed(props, myself.Name)

	if err != nil {
		println("[PEER] Name already in use")
	}

	console.ReadLine()

}
