package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/mlisa/gomath/common"
	"github.com/mlisa/gomath/message"
)

type Controller struct {
	Gui            *GuiCoordinator
	CoordinatorPID *actor.PID
	Coordinator    *Coordinator
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
	remote.Start(config.Address)
	list, err := common.GetCoordinatorsList()
	if err != nil {
		return err
	}
	c.Coordinator = &Coordinator{MaxPeers: *maxpeer, Peers: make(map[string]*actor.PID), Coordinators: list, Controller: c}
	props := actor.FromInstance(c.Coordinator)
	c.CoordinatorPID, err = actor.SpawnNamed(props, config.Id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) RunPing() {
	if c.CoordinatorPID != nil {
		c.CoordinatorPID.Tell(&message.Ping{})
	}
}

func (c *Controller) UpdatePings(pings map[string]common.Pong) {
	c.Gui.UpdatePings(pings)
}

func (c *Controller) Log(s string) {
	c.Gui.PrintToView("log", fmt.Sprintf("[%s] %s", c.CoordinatorPID.Id, s))
}

func (c *Controller) GetLatency(peer string) int64 {
	req := c.CoordinatorPID.RequestFuture(&message.GetPing{Peer: peer}, 1*time.Second)
	if r, err := req.Result(); err == nil {
		return r.(*message.Pong).Pong
	}
	return -1
}

func (c *Controller) GetPeers() map[string]*actor.PID {
	return c.Coordinator.Peers
}
