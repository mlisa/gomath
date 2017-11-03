package main

import (
	"log"
	"strconv"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/fatih/color"
)

type PublishActor struct {
	AppPath   string
	RabbitPid string
}

func (state *PublishActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		color.Cyan("Started, initialize actor here")
		context.SetBehavior(running)
	case *actor.Stopping:
		color.Red("Stopping, actor is about shut down")
	case *actor.Stopped:
		color.Green("Stopped, actor and it's children are stopped")
	case *actor.Restarting:
		color.Magenta("Restarting, actor is about restart")
	default:
		log.Println(msg)
	}
}

func running(context actor.Context) {
	switch msg := context.Message().(type) {
	case int:
		color.Yellow("INT: " + strconv.Itoa(msg))
	case string:
		color.Yellow("STRING:" + msg)
	}
}

func main() {
	remote.Start("localhost:0")
	props := actor.FromInstance(&PublishActor{})
	publishPid := actor.Spawn(props)

	publishPid.Tell(12)
	publishPid.Tell("ASDSADDSA")
	time.Sleep(2 * time.Second)
	publishPid.Tell("ASDSADDSA")

	console.ReadLine()
}
