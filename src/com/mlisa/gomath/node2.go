package main

import (
	"com/mlisa/gomath/message"
	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"log"
	"runtime"
	//"github.com/AsynkronIT/protoactor-go/examples/distributedchannels/messages"
)
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	remote.Start("127.0.0.1:8080")
	//create the channel
	//channel := make(chan *message.Hello)

	//create an actor receiving messages and pushing them onto the channel
	props := actor.FromFunc(func(context actor.Context) {
		if msg, ok := context.Message().(*message.Hello); ok {
			log.Println(msg)
			//channel <- msg
		}
	})


	actor.SpawnNamed(props, "MyMessage")

	//consume the channel just like you use to
	/*go func() {
		for msg := range channel {
			//log.Println(msg)
		}
	}()*/

	console.ReadLine()
}
