package main

import (
	"bytes"
	"com/mlisa/gomath/message"
	"com/mlisa/gomath/parser"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/jroimartin/gocui"
	"log"
)

type Controller struct {
	gui   *gocui.Gui
	peer  *actor.PID
	cache *CacheManager
}

func (controller *Controller) Receive(context actor.Context) {

	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Println("[CONTROLLER] Started, initialize actor here, I'm " + context.Self().Id + " " + context.Self().Address)
		controller.peer.Request(message.Hello{}, context.Self())

	case *message.AskForResult:
		//Policy per decidere se calcolarlo qui o meno
		var something = true
		if something {
			controller.peer.Request(msg, context.Self())
			controller.gui.Update(func(g *gocui.Gui) error {
				output, _ := g.View("Log")
				fmt.Fprint(output, "Sent AskForResult message to peer")
				return nil
			})
		} else {
			result, err := parser.ParseReader("", bytes.NewBufferString(msg.Operation))
			if err == nil {
				controller.gui.Update(func(g *gocui.Gui) error {
					output, _ := g.View("Output")
					fmt.Fprint(output, result)
					return nil
				})
				controller.cache.addNewOperation(msg.Operation, string(result.(int)))
			}
		}

	case *message.Response:
		controller.gui.Update(func(g *gocui.Gui) error {
			output, _ := g.View("Output")
			fmt.Fprint(output, msg.Result)
			return nil
		})
		controller.gui.Update(func(g *gocui.Gui) error {
			output, _ := g.View("Log")
			fmt.Fprint(output, "Receive Response message from peer")
			return nil
		})

	case *message.SearchInCache:
		result, found := controller.cache.retrieveResult(msg.Operation)
		controller.gui.Update(func(g *gocui.Gui) error {
			output, _ := g.View("Log")
			fmt.Fprint(output, "Receive SearchInCache message from peer")
			return nil
		})

		if found {
			context.Sender().Tell(message.ResponseFromCache{result, msg.FromPeer})
		}

	}
}
