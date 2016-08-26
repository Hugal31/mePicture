package config

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/BurntSushi/toml"
)

type Config struct {
	DatabaseFile string
	PicturesRoot string
}

var config *Config = readConfig()

func GetConfig() *Config {
	return config
}

func getHomePath() string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return currentUser.HomeDir
}

func getConfigFileName() string {
	return getHomePath() + "/.mePicture"
}

func readConfig() *Config {
	var config Config
	if _, err := toml.DecodeFile(getConfigFileName(), &config); err != nil {
		homePath := getHomePath()
		config.DatabaseFile = homePath + "/.mePicture.sql"
		config.PicturesRoot = homePath + "/Pictures/Wallpapers"
		fmt.Fprintf(os.Stderr, "Assuming pictures are in %s\n", config.PicturesRoot)
	}
	return &config
}

func SaveConfig() {
}
