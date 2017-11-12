package common

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type Node struct {
	Id                    string
	Address               string
	Latency               float32
	ComputationCapability float32
	Queue                 float32
}

type Config struct {
	Myself       Node
	Coordinators []*actor.PID
}

type jsonParse struct {
	Id      string `json:"id"`
	Address string `json:"address"`
}

func GetFileConfig(path string) (Config, error) {
	fileAbs, _ := filepath.Abs(path)
	configuration := Config{}

	if _, err := os.Stat(fileAbs); !os.IsNotExist(err) {
		file, err := os.Open(fileAbs)
		if err != nil {
			return configuration, err
		}
		defer file.Close()
		if err = json.NewDecoder(file).Decode(&configuration); err != nil {
			return configuration, err
		}
		return configuration, nil
	} else {
		return configuration, err
	}
}

func GetCoordinatorsList() (map[string]*actor.PID, error) {
	url := "http://gomath.duckdns.org:8080/mirror.json"
	client := &http.Client{Timeout: 10 * time.Second}
	var list []jsonParse
	coordinators := make(map[string]*actor.PID, len(list))

	r, err := client.Get(url)
	if err != nil {
		return coordinators, err
	}
	defer r.Body.Close()
	if r != nil && err == nil {
		// read []byte{}
		b, _ := ioutil.ReadAll(r.Body)

		// Due to some presence of unicode chars convert raw JSON to string than parse it
		// GO strings works with utf-8
		if err = json.NewDecoder(strings.NewReader(string(b))).Decode(&list); err != nil {
			return coordinators, err
		}
	}
	for i := range list {
		coordinators[list[i].Id] = actor.NewPID(list[i].Address, list[i].Id)
	}
	return coordinators, nil
}
