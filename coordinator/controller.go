package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/jroimartin/gocui"
	"github.com/mlisa/gomath/common"
	"github.com/mlisa/gomath/message"
)

type Controller struct {
	Gui         *gocui.Gui
	Coordinator *actor.PID
}

func (c *Controller) PublishCoordinator(token string) {
	url := "http://gomath.duckdns.org:8080/publish"

	jsonStr := map[string]string{"token": token}
	jsonReq, _ := json.Marshal(jsonStr)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
}

func (c *Controller) StartCoordinator(config common.Config) error {
	remote.Start(config.Myself.Address)
	list, err := common.GetCoordinatorsList()
	if err != nil {
		return err
	}
	props := actor.FromInstance(&Coordinator{MaxPeers: *maxpeer, Peers: make(map[string]*actor.PID), Coordinators: list, Controller: c})
	_, err = actor.SpawnNamed(props, config.Myself.Id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) RunPing() {
	if c.Coordinator != nil {
		c.Coordinator.Tell(&message.Ping{})
	}
}
