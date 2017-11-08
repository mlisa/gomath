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

func GetConfig(who string) Config {
	fileName := "config_" + who + ".json"
	absPath, _ := filepath.Abs(fileName)
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
