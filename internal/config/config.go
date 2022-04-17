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
func (c *Config) Load(configPath string) (err error) {
	// Get the path to the config file
	configPath = c.GetConfigPath(configPath)

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
	err = jsonParser.Decode(&c)
	if err != nil {
		return err
	}

	return nil
}

// Save saves a SoftTube configuration file
func (c *Config) Save(configPath string) {
	// Get the path to the config file
	configPath = c.GetConfigPath(configPath)

	// Open config file
	configFile, err := os.OpenFile(configPath, os.O_TRUNC|os.O_WRONLY, 0644)

	// Handle errors
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Create JSON from config object
	data, err := json.MarshalIndent(c, "", "\t")

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
func (c *Config) GetConfigPath(configPath string) string {
	if configPath == "" {
		configPath = defaultConfigFile
	}

	home := c.getHomeDirectory()
	return path.Join(home, configPath)
}

// ClearRepositories clears the slice of repositories
func (c *Config) ClearRepositories() {
	c.Repositories = nil
}

// AddRepository adds a new repository
func (c *Config) AddRepository(path, imagePath string, isFavorite bool) {
	repo := &Repository{Path: path, ImagePath: imagePath, IsFavorite: isFavorite}
	c.Repositories = append(c.Repositories, repo)
}

// RemoveRepository adds a new repository
func (c *Config) RemoveRepository(path string) {
	for i := range c.Repositories {
		if c.Repositories[i].Path == path {
			c.Repositories = append(c.Repositories[:i], c.Repositories[i+1:]...)
		}
	}
}

// GetExternalApplicationByName gets an external application by name
func (c *Config) GetExternalApplicationByName(name string) *ExternalApplication {
	for i := range c.ExternalApplications {
		ext := c.ExternalApplications[i]
		if ext.Name == name {
			return ext
		}
	}

	return nil
}

// ClearExternalApplications clears the slice of external applications
func (c *Config) ClearExternalApplications() {
	c.ExternalApplications = nil
}

// AddExternalApplication adds an external application
func (c *Config) AddExternalApplication(name, command, argument string) {
	a := &ExternalApplication{
		Name:     name,
		Command:  command,
		Argument: argument,
	}

	c.ExternalApplications = append(c.ExternalApplications, a)
}

// RemoveExternalApplication adds a new extenal application
func (c *Config) RemoveExternalApplication(name string) {
	for i := range c.ExternalApplications {
		if c.ExternalApplications[i].Name == name {
			c.ExternalApplications = append(
				c.ExternalApplications[:i],
				c.ExternalApplications[i+1:]...,
			)
		}
	}
}

//
// Private functions
//

// Get current users home directory
func (c *Config) getHomeDirectory() string {
	u, err := user.Current()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get user home directory : %s", err)
		panic(errorMessage)
	}
	return u.HomeDir
}
