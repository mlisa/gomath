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
	Gui       *gocui.Gui
	Peer      *actor.PID
	Cache     *CacheManager
	Config    common.Config
	Connected bool
	Latency   int64
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
	ASKCOORDINATOR
	EXTERNALANSWER
	UNABLETOCONNECT
)

func (controller *Controller) AskForResult(operation string) {
	operation = strings.TrimSpace(operation)
	if controller.Connected {
		if _, err := parser.Parse("", []byte(operation)); err == nil {
			if resultInLocalCache := controller.SearchInCache(operation); resultInLocalCache != "" {
				controller.SetOutput(resultInLocalCache)
				controller.Log(FOUNDRESULTINCACHE, "")
			} else {
				var complexity = int64(strings.Count(operation, "*")*2 + strings.Count(operation, "/")*2 +
					strings.Count(operation, "+") + strings.Count(operation, "-"))

				if complexity*100 > controller.Config.ComputeCapability || controller.Latency > 450 {
					controller.Peer.Tell(&message.AskForResult{operation})
					controller.Log(ASKFORRESULT, "")
				} else {
					controller.ComputeLocal(operation)
				}
			}
		} else {
			controller.SetOutput("[ERROR] Wrong input format")
		}
	} else {
		controller.ComputeLocal(operation)
	}

}

func (controller *Controller) ComputeLocal(operation string) {
	result, _ := parser.ParseReader("", bytes.NewBufferString(operation))

	controller.SetOutput(strconv.Itoa(result.(int)))
	controller.Log(OFFLINECOMPUTATION, "")
	controller.Cache.addNewOperation(operation, strconv.Itoa(result.(int)))

}

func (controller *Controller) SearchInCache(operation string) string {
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

func (controller *Controller) Log(eventType EventType, from string) {
	switch eventType {
	case NEWNODE:
		controller.setLog("[" + controller.Peer.Id + "] New node entered in region: " + from)

	case DEADNODE:
		controller.setLog("[" + controller.Peer.Id + "] One node of the region died: " + from)

	case LOSTCONNECTION:
		controller.setLog("[" + controller.Peer.Id + "] Lost connection from coordinator: " + from)

	case FOUNDNEWCOORDINATOR:
		controller.setLog("[" + controller.Peer.Id + "] Found new coordinator: " + from)

	case ASKFORRESULT:
		controller.setLog("[" + controller.Peer.Id + "] Sent AskForResult message to peers")

	case SEARCHINCACHE:
		controller.setLog("[" + controller.Peer.Id + "] Received RequestForCache message from: " + from)

	case RECEIVEDRESPONSE:
		controller.setLog("[" + controller.Peer.Id + "] Received Response message from " + from)

	case FOUNDRESULTINCACHE:
		controller.setLog("[" + controller.Peer.Id + "] Retrieved result from local cache")

	case OFFLINECOMPUTATION:
		controller.setLog("[" + controller.Peer.Id + "] Operation computed offline")

	case NOTFOUND:
		controller.setLog("[" + controller.Peer.Id + "] Operation not present in cache")

	case NORESPONSE:
		controller.setLog("[" + controller.Peer.Id + "] No one has the answer... ")

	case EXTERNALANSWER:
		controller.setLog("[" + controller.Peer.Id + "] Received answer from another region")

	case ASKCOORDINATOR:
		controller.setLog("[" + controller.Peer.Id + "] No one has the response, contacting coordinator...")

	case UNABLETOCONNECT:
		controller.setLog("[" + controller.Peer.Id + "] Unable to connect to GoMath system")
	}
}
