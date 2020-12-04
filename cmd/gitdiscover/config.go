package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"
)

type Config struct {
	Paths           []string `json:"paths"`
	DateFormat      string   `json:"date-format"`
	PathColumnWidth int      `json:"path-column-width"`
}

func NewConfig() *Config {
	return new(Config)
}

func (config *Config) ConfigExists() bool {
	// Get the path to the config file
	configPath := getConfigPath()

	_, err := os.Stat(configPath)
	return !os.IsNotExist(err)
}

// Load : Loads the configuration file
func (config *Config) Load() (err error) {
	// Get the path to the config file
	configPath := getConfigPath()

	// Open config file
	configFile, err := os.Open(configPath)

	// Handle errors
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	defer func() {
		err = configFile.Close()
	}()

	// Parse the JSON document
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err!=nil {
		return err
	}

	return nil
}

// Get path to the config file
func getConfigPath() string {
	home := getHomeDirectory()

	return path.Join(home, ".config/softteam/gitdiscover/config.json")
}

// Get current users home directory
func getHomeDirectory() string {
	u, err := user.Current()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get user home directory : %s", err)
		panic(errorMessage)
	}
	return u.HomeDir
}
