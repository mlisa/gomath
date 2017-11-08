package main

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/mlisa/gomath/message"
	"github.com/mlisa/gomath/parser"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/jroimartin/gocui"
)

type Controller struct {
	Gui   *gocui.Gui
	Peer  *actor.PID
	Cache *CacheManager
}

func (controller *Controller) AskForResult(operation string) {
	var something = true
	if something {
		controller.Peer.Tell(&message.AskForResult{operation})
		controller.Gui.Update(func(g *gocui.Gui) error {
			output, _ := g.View("Log")
			output.Clear()
			fmt.Fprintln(output, "Sent AskForResult message to peer")
			return nil
		})
	} else {
		result, err := parser.ParseReader("", bytes.NewBufferString(operation))
		if err == nil {
			controller.Gui.Update(func(g *gocui.Gui) error {
				output, _ := g.View("Output")
				output.Clear()
				fmt.Fprint(output, result)
				return nil
			})
			controller.Gui.Update(func(g *gocui.Gui) error {
				output, _ := g.View("Log")
				output.Clear()
				fmt.Fprintln(output, "Operation computed offline")
				return nil
			})
			controller.Cache.addNewOperation(operation, strconv.Itoa(result.(int)))
		}
	}
}

func (controller *Controller) SearchInCache(operation string) string {
	controller.Gui.Update(func(g *gocui.Gui) error {
		output, _ := g.View("Log")
		fmt.Fprintln(output, "Received SearchInCache message from peer")
		return nil
	})

	if result, err := controller.Cache.retrieveResult(operation); err == nil {
		controller.Gui.Update(func(g *gocui.Gui) error {
			output, _ := g.View("Log")
			fmt.Fprintln(output, "Retrieved result from local cache")
			return nil
		})
		return result
	}
	return ""
}

func (controller *Controller) SetResult(result string) {
	controller.Gui.Update(func(g *gocui.Gui) error {
		output, _ := g.View("Output")
		fmt.Fprintln(output, result)
		return nil
	})
	controller.Gui.Update(func(g *gocui.Gui) error {
		output, _ := g.View("Log")
		fmt.Fprintln(output, "Receive Response message from peer")
		return nil
	})
}
