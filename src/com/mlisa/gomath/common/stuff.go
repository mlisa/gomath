package common

import (
	"encoding/json"
	//"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/actor"
	"log"
	"os"
	"path/filepath"
)

type Peer struct {
	Id      string
	Address string
}

type Config struct {
	Myself       Peer
	Coordinators []*actor.PID
}

func GetConfig(who string) Config {
	fileName := "config_" + who + ".json"
	absPath, _ := filepath.Abs(filepath.Clean(fileName))
	configuration := Config{}
	file, err := os.Open(absPath)
	if err != nil {
		log.Println("[ERROR] " + err.Error())
	}
	defer file.Close()
	if err = json.NewDecoder(file).Decode(&configuration); err != nil {
		log.Fatalln("[ERROR] " + err.Error())
	}
	return configuration
}
