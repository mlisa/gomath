package Common

import (
	"path/filepath"
	"os"
	"encoding/json"
	"log"
)

type PID struct {
	Name string
	Address string
}

type Config struct {
	Myself PID
	Coordinators []PID
}

func getConfig() Config {
	absPath, _ := filepath.Abs("com/mlisa/gomath/config.json")
	file, err := os.Open(absPath)
	if err != nil {
		log.Println("[ERROR] " + err.Error())
	}
	decoder := json.NewDecoder(file)
	configuration := Config{}
	decoder.Decode(&configuration)
	return configuration
}