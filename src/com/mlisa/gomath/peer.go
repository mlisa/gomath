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

	switch context.Message().(type) {
	case *actor.Started:
		fmt.Println("[PEER] Started, initialize actor here, I'm " + context.Self().Id + " " + context.Self().Address)

		coordinators := getConfig().Coordinators //lettura da file config
		for _, PID := range coordinators {
			log.Println("[PEER] Try to connect to " + PID.Address + " " + PID.Name)
			coordinator := actor.NewPID(PID.Address, PID.Name)
			coordinator.Request(&message.Hello{context.Self()}, context.Self())
		}

	case *message.Available:
		log.Println("[PEER] Found a coordinator!")
		coord = context.Sender()
		coord.Tell(&message.Register{context.Self()})
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
			response := &message.Response{context.Self(), result}
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
