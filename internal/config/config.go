package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"
)

// Config : The main config type
type Config struct {
	Repositories         []*Repository          `json:"repositories"`
	ExternalApplications []*ExternalApplication `json:"external-applications"`
	DateFormat           string                 `json:"date-format"`
	PathColumnWidth      int                    `json:"path-column-width"`
}

// Repository : A Repository in the config
type Repository struct {
	Path       string `json:"path"`
	ImagePath  string `json:"image-path"`
	IsFavorite bool   `json:"is-favorite"`
}

// ExternalApplication : An external application in the config
type ExternalApplication struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	Argument string `json:"argument"`
}

// NewConfig creates a new config
func NewConfig() *Config {
	return new(Config)
}

// Load loads the configuration file
func (config *Config) Load(configPath string) (err error) {
	// Get the path to the config file
	configPath = config.GetConfigPath(configPath)

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
	if err != nil {
		return err
	}

	return nil
}

// Save saves a SoftTube configuration file
func (config *Config) Save(configPath string) {
	// Get the path to the config file
	configPath = config.GetConfigPath(configPath)

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
	_, _ = configFile.Write(data)

	_ = configFile.Close()
}

// GetConfigPath returns the path to the config file
func (config *Config) GetConfigPath(configPath string) string {
	if configPath == "" {
		configPath = defaultConfigFile
	}

	home := config.getHomeDirectory()
	return path.Join(home, configPath)
}

// ClearRepositories clears the slice of repositories
func (config *Config) ClearRepositories() {
	config.Repositories = nil
}

// AddRepository adds a new repository
func (config *Config) AddRepository(path, imagePath string) {
	repo := &Repository{Path: path, ImagePath: imagePath}
	config.Repositories = append(config.Repositories, repo)
}

// RemoveRepository adds a new repository
func (config *Config) RemoveRepository(path string) {
	for i := range config.Repositories {
		if config.Repositories[i].Path == path {
			config.Repositories = append(config.Repositories[:i], config.Repositories[i+1:]...)
		}
	}
}

// GetExternalApplicationByName gets an external application by name
func (config *Config) GetExternalApplicationByName(name string) *ExternalApplication {
	for i := range config.ExternalApplications {
		ext := config.ExternalApplications[i]
		if ext.Name == name {
			return ext
		}
	}

	return nil
}

// ClearExternalApplications clears the slice of external applications
func (config *Config) ClearExternalApplications() {
	config.ExternalApplications = nil
}

// AddExternalApplication adds an external application
func (config *Config) AddExternalApplication(name, command, argument string) {
	a := &ExternalApplication{
		Name:     name,
		Command:  command,
		Argument: argument,
	}

	config.ExternalApplications = append(config.ExternalApplications, a)
}

// RemoveExternalApplication adds a new extenal application
func (config *Config) RemoveExternalApplication(name string) {
	for i := range config.ExternalApplications {
		if config.ExternalApplications[i].Name == name {
			config.ExternalApplications = append(
				config.ExternalApplications[:i],
				config.ExternalApplications[i+1:]...,
			)
		}
	}
}

//
// Private functions
//

// Get current users home directory
func (config *Config) getHomeDirectory() string {
	u, err := user.Current()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get user home directory : %s", err)
		panic(errorMessage)
	}
	return u.HomeDir
}
