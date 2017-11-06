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
	Peer  *actor.PID
	cache *CacheManager
}

func (controller *Controller) AskForResult(operation string) {
	var something = true
	if something {
		controller.Peer.Tell(&message.AskForResult{operation})
		controller.gui.Update(func(g *gocui.Gui) error {
			output, _ := g.View("Log")
			fmt.Fprint(output, "Sent AskForResult message to peer")
			return nil
		})
	} else {
		result, err := parser.ParseReader("", bytes.NewBufferString(operation))
		if err == nil {
			controller.gui.Update(func(g *gocui.Gui) error {
				output, _ := g.View("Output")
				fmt.Fprint(output, result)
				return nil
			})
			log.Println("[CONTROLLER] Saved result :" + strconv.Itoa(result.(int)))
			controller.cache.addNewOperation(operation, strconv.Itoa(result.(int)))
		}
	}
}

func (controller *Controller) SearchInCache(operation string) string {
	controller.gui.Update(func(g *gocui.Gui) error {
		output, _ := g.View("Log")
		fmt.Fprint(output, "Received SearchInCache message from peer")
		return nil
	})

	result, found := controller.cache.retrieveResult(operation)
	if found {
		log.Println("Retrieved result " + result + " from cache")
		return result
	}
	return ""
}

func (controller *Controller) SetResult(result string) {
	log.Println("Received result " + result)
	controller.gui.Update(func(g *gocui.Gui) error {
		output, _ := g.View("Output")
		fmt.Fprint(output, result)
		return nil
	})
	controller.gui.Update(func(g *gocui.Gui) error {
		output, _ := g.View("Log")
		fmt.Fprint(output, "Receive Response message from peer")
		return nil
	})
}
