package main

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/mlisa/gomath/message"
	"github.com/mlisa/gomath/parser"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

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
)

type Coordinators struct {
	Id      string `json:"id"`
	Address string `json:"address"`
}

func (c *Controller) getCoordinatorsList() ([]Coordinators, error) {
	url := "http://gomath.duckdns.org:8080/mirror.json"
	client := &http.Client{Timeout: 10 * time.Second}
	var out []Coordinators

	r, err := client.Get(url)
	if err != nil {
		return out, err
	}
	defer r.Body.Close()
	if r != nil && err == nil {
		// read []byte{}
		b, _ := ioutil.ReadAll(r.Body)

		// Due to some presence of unicode chars convert raw JSON to string than parse it
		// GO strings works with utf-8
		if err = json.NewDecoder(strings.NewReader(string(b))).Decode(&out); err != nil {
			return out, err
		}
	}
	return out, nil

}

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

	}
}
