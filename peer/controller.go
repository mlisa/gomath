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
		controller.SetOutput(result, err)
		if err == nil {
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

func (controller *Controller) SetOutput(result interface{}, err error) {
	outputString := ""
	if err != nil {
		outputString = "[ERROR] Wrong input format"
	} else {
		outputString = result.(string)
	}
	controller.Gui.Update(func(g *gocui.Gui) error {
		output, _ := g.View("Output")
		output.Clear()
		fmt.Fprintln(output, outputString)
		return nil
	})
}

func (controller *Controller) setLog(log string) {
	controller.Gui.Update(func(g *gocui.Gui) error {
		output, _ := g.View("Log")
		fmt.Fprintln(output, log)
		return nil
	})
}
func (controller *Controller) Log(eventType EventType) {
	switch eventType {
	case NEWNODE:
		controller.setLog("New node entered in region")

	case DEADNODE:
		controller.setLog("One node of the region died")

	case LOSTCONNECTION:
		controller.setLog("Lost connection from coordinator")

	case FOUNDNEWCOORDINATOR:
		controller.setLog("Found new coordinator.")

	case ASKFORRESULT:
		controller.setLog("Sent AskForResult message to peers")

	case SEARCHINCACHE:
		controller.setLog("Received SearchInCache message from peer")

	case RECEIVEDRESPONSE:
		controller.setLog("Received Response message from peer")

	case FOUNDRESULTINCACHE:
		controller.setLog("Retrieved result from local cache")

	case OFFLINECOMPUTATION:
		controller.setLog("Operation computed offline")

	}
}
