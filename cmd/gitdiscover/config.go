package main

import (
"encoding/json"
"fmt"
"os"
"os/user"
"path"
)

type Config struct {
	Paths []string `json:"paths"`
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
func (config *Config) Load() error {
	// Get the path to the config file
	configPath := getConfigPath()

	// Open config file
	configFile, err := os.Open(configPath)

	// Handle errors
	if err != nil {
		fmt.Println(err.Error())
	}
	defer configFile.Close()

	// Parse the JSON document
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return nil
}


// Get path to the config file
func getConfigPath() string {
	home := getHomeDirectory()

	return path.Join(home, ".config/gitdiscover/config.json")
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

