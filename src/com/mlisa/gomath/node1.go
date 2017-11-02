package main

/*import (
	//"fmt"
	"com/mlisa/gomath/message"
	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"runtime"
	//"github.com/AsynkronIT/protoactor-go/examples/distributedchannels/messages"
	//"fmt"
	"github.com/AsynkronIT/protoactor-go/remote"
)

func newMyMessageSenderChannel() chan<- *message.Hello {
	channel := make(chan *message.Hello)
	remote := actor.NewPID("127.0.0.1:8080", "MyMessage")
	go func() {
		message := &message.Hello{
			Content: "hello",
		}
		for msg := range channel {
			remote.Tell(msg)
		}
	}()

	return channel
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	remote.Start("127.0.0.1:0")

	for i := 0; i < 10; i++ {
		message := &message.Hello{
			Content: "hello",
		}
		channel <- message
	}

	console.ReadLine()
}
*/