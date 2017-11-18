package main

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/mlisa/gomath/message"
	"github.com/mlisa/gomath/parser"

	"strings"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/jroimartin/gocui"
	"github.com/mlisa/gomath/common"
)

type Controller struct {
	Gui    *gocui.Gui
	Peer   *actor.PID
	Cache  *CacheManager
	Config common.Config
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
	NOTFOUND
	NORESPONSE
)

func (controller *Controller) AskForResult(operation string) {
	operation = strings.TrimSpace(operation)
	var complexity = strings.Count(operation, "*")*2 + strings.Count(operation, "/")*2 +
		strings.Count(operation, "+") + strings.Count(operation, "-")
	if float32(complexity*100) > controller.Config.Myself.ComputationCapability {
		controller.Peer.Tell(&message.AskForResult{operation})
		controller.Log(ASKFORRESULT)
	} else {
		controller.ComputeLocal(operation)
	}
}

func (controller *Controller) ComputeLocal(operation string) {
	result, err := parser.ParseReader("", bytes.NewBufferString(operation))
	if err == nil {
		controller.SetOutput(strconv.Itoa(result.(int)))
		controller.Log(OFFLINECOMPUTATION)
		controller.Cache.addNewOperation(operation, strconv.Itoa(result.(int)))
	} else {
		controller.SetOutput("[ERROR] Wrong input format")
	}
}

func (controller *Controller) SearchInCache(operation string) string {
	controller.Log(SEARCHINCACHE)
	if result, err := controller.Cache.retrieveResult(operation); err == nil {
		return result
	}
	return ""
}

func (controller *Controller) SetOutput(outputString string) {
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
		controller.setLog("Received RequestForCache message from peer")

	case RECEIVEDRESPONSE:
		controller.setLog("Received Response message from peer")

	case FOUNDRESULTINCACHE:
		controller.setLog("Retrieved result from local cache")

	case OFFLINECOMPUTATION:
		controller.setLog("Operation computed offline")

	case NOTFOUND:
		controller.setLog("Operation not present in cache")

	case NORESPONSE:
		controller.setLog("No one has the answer... ")

	}
}
