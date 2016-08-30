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

var config *Config

func GetConfig() *Config {
	if config == nil {
		config = readConfig()
	}
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
		config.DatabaseFile = homePath + string(os.PathSeparator) + ".mePicture.sql"
		config.PicturesRoot = string(os.PathSeparator)
		config.Save()
		fmt.Fprintf(os.Stderr, "Picture root set to %s\n", config.PicturesRoot)
	}
	return &config
}

func (conf *Config) Save() {
	file, err := os.OpenFile(getConfigFileName(), os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	encoder := toml.NewEncoder(file)
	fmt.Print(encoder)
	encoder.Encode(conf)
}
