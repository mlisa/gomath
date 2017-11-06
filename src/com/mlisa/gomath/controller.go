package main

import (
	"bytes"
	"com/mlisa/gomath/message"
	"com/mlisa/gomath/parser"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/jroimartin/gocui"
	"log"
	"strconv"
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
		log.Println("[CONTROLLER] Sending hello to " + controller.peer.Id + " " + controller.peer.Address)
		controller.peer.Request(&message.Hello{}, context.Self())
		//context.Request(controller.peer, message.Hello{})

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
				log.Println("[CONTROLLER] Saved result :" + strconv.Itoa(result.(int)))
				controller.cache.addNewOperation(msg.Operation, strconv.Itoa(result.(int)))
			}
		}

	case *message.Response:
		log.Println("Received result " + msg.Result)
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
		controller.gui.Update(func(g *gocui.Gui) error {
			output, _ := g.View("Log")
			fmt.Fprint(output, "Received SearchInCache message from peer")
			return nil
		})

		result, found := controller.cache.retrieveResult(msg.Operation)
		if found {
			log.Println("Retrieved result " + result + " from cache")
			context.Sender().Tell(&message.ResponseFromCache{result, msg.FromPeer})
		}

	}
}
