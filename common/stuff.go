package common

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type Peer struct {
	Id      string
	Address string
}

type Config struct {
	Myself       Peer
	Coordinators []*actor.PID
}

func GetFileConfig(path string) Config {
	absPath, _ := filepath.Abs(filepath.Clean(path))
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
