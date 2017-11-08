package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/jroimartin/gocui"
	"github.com/mlisa/gomath/common"
)

type Controller struct {
	Gui         *gocui.Gui
	Coordinator *actor.PID
}

type Coordinators struct {
	Id      string `json:"id"`
	Address string `json:"address"`
}

func (c *Controller) getCoordinatorsList() ([]Coordinators, error) {
	url := "http://172.17.0.2/mirror.json"
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

func (c *Controller) PublishCoordinator(token string) {
	url := "http://172.17.0.2/publish"

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
	list, err := c.getCoordinatorsList()
	if err != nil {
		return err
	}
	props := actor.FromInstance(&Coordinator{MaxPeers: *maxpeer, Peers: make(map[string]*actor.PID), Coordinators: list})
	_, err = actor.SpawnNamed(props, config.Myself.Id)
	if err != nil {
		return err
	}
	return nil
}
