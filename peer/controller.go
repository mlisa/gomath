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

type EventType int

const (
	NEWNODE EventType = 1 + iota
	DEADNODE
	LOSTCONNECTION
	FOUNDNEWCOORDINATOR
	ASKFORRESULT
	SEARCHINCACHE
	RECEIVEDRESPONSE
	OFFLINECOMPUTATION
	FOUNDRESULTINCACHE
)

func (controller *Controller) AskForResult(operation string) {
	var something = true
	if something {
		controller.Peer.Tell(&message.AskForResult{operation})
		controller.Log(ASKFORRESULT)
	} else {
		result, err := parser.ParseReader("", bytes.NewBufferString(operation))
		if err == nil {
			controller.Gui.Update(func(g *gocui.Gui) error {
				output, _ := g.View("Output")
				output.Clear()
				fmt.Fprint(output, result)
				return nil
			})
			controller.Log(OFFLINECOMPUTATION)
			controller.Cache.addNewOperation(operation, strconv.Itoa(result.(int)))
		}
	}
}

func (controller *Controller) SearchInCache(operation string) string {
	controller.Log(SEARCHINCACHE)

	if result, err := controller.Cache.retrieveResult(operation); err == nil {
		controller.Log(FOUNDRESULTINCACHE)
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
	controller.Log(RECEIVEDRESPONSE)
}

func (controller *Controller) Log(eventType EventType) {
	switch eventType {
	case NEWNODE:
		controller.Gui.Update(func(g *gocui.Gui) error {
			output, _ := g.View("Log")
			fmt.Fprintln(output, "New node entered in region")
			return nil
		})
	case DEADNODE:
		controller.Gui.Update(func(g *gocui.Gui) error {
			output, _ := g.View("Log")
			fmt.Fprintln(output, "One node of the region died")
			return nil
		})
	case LOSTCONNECTION:
		controller.Gui.Update(func(g *gocui.Gui) error {
			output, _ := g.View("Log")
			fmt.Fprintln(output, "Lost connection from coordinator")
			return nil
		})
	case FOUNDNEWCOORDINATOR:
		controller.Gui.Update(func(g *gocui.Gui) error {
			output, _ := g.View("Log")
			fmt.Fprintln(output, "Found new coordinator.")
			return nil
		})
	case ASKFORRESULT:
		controller.Gui.Update(func(g *gocui.Gui) error {
			output, _ := g.View("Log")
			fmt.Fprintln(output, "Sent AskForResult message to peers")
			return nil
		})
	case SEARCHINCACHE:
		controller.Gui.Update(func(g *gocui.Gui) error {
			output, _ := g.View("Log")
			fmt.Fprintln(output, "Received SearchInCache message from peer")
			return nil
		})
	case RECEIVEDRESPONSE:
		controller.Gui.Update(func(g *gocui.Gui) error {
			output, _ := g.View("Log")
			fmt.Fprintln(output, "Receive Response message from peer")
			return nil
		})
	case FOUNDRESULTINCACHE:
		controller.Gui.Update(func(g *gocui.Gui) error {
			output, _ := g.View("Log")
			fmt.Fprintln(output, "Retrieved result from local cache")
			return nil
		})
	case OFFLINECOMPUTATION:
		controller.Gui.Update(func(g *gocui.Gui) error {
			output, _ := g.View("Log")
			fmt.Fprintln(output, "Operation computed offline")
			return nil
		})

	}
}
