package common

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type PID struct {
	Name    string
	Address string
}

type Config struct {
	Myself       PID
	Coordinators []PID
}

func getConfig() Config {
	absPath, _ := filepath.Abs("com/mlisa/gomath/config.json")
	file, err := os.Open(absPath)
	if err != nil {
		log.Println("[ERROR] " + err.Error())
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Config{}
	decoder.Decode(&configuration)
	return configuration
}
