package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"
)

const emptyConfig = `{
  "paths": [
  ],
  "date-format": "2006-01-02, kl. 15:04",
  "path-column-width":40
}`

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
	configPath := config.GetConfigPath()

	_, err := os.Stat(configPath)
	return !os.IsNotExist(err)
}

func (config *Config) CreateEmptyConfig() error {
	f, err := os.Create(config.GetConfigPath())
	if err != nil {
		return err
	}
	_, err = f.WriteString(emptyConfig)
	if err != nil {
		_ = f.Close()
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

// Load : Loads the configuration file
func (config *Config) Load() (err error) {
	// Get the path to the config file
	configPath := config.GetConfigPath()

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

// Save : Saves a SoftTube configuration file
func (config *Config) Save() {
	// Get the path to the config file
	configPath := config.GetConfigPath()

	// Open config file
	configFile, err := os.OpenFile(configPath, os.O_TRUNC|os.O_WRONLY, 0644)

	// Handle errors
	if err != nil {
		fmt.Println(err.Error())
		return
	}


	// Create JSON from config object
	data, err := json.MarshalIndent(config, "", "\t")

	// Handle errors
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Write the data
	_,_ = configFile.Write(data)

	_ = configFile.Close()
}

// Get path to the config file
func (config *Config) GetConfigPath() string {
	home := config.getHomeDirectory()

	return path.Join(home, ".config/softteam/gitdiscover/config.json")
}

// Get current users home directory
func (config *Config) getHomeDirectory() string {
	u, err := user.Current()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get user home directory : %s", err)
		panic(errorMessage)
	}
	return u.HomeDir
}
