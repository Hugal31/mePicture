package config

import (
	"os/user"
	"log"
	"github.com/BurntSushi/toml"
)

type Config struct {
	DatabaseFile string
}

var config *Config = readConfig()

func GetConfig() *Config {
	return config
}

func getConfigFileName() string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return currentUser.HomeDir + "/.mePicture"
}

func readConfig() *Config {
	var config Config
	if _, err := toml.DecodeFile(getConfigFileName(), &config); err != nil {
		currentUser, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		config.DatabaseFile = currentUser.HomeDir + "/.mePicture.sql"
	}
	return &config
}

func SaveConfig() {
}
