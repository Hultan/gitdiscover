package gitdiscover

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"
)

// Config : The main config type
type Config struct {
	Repositories         []Repository          `json:"repositories"`
	ExternalApplications []ExternalApplication `json:"external-applications"`
	DateFormat           string                `json:"date-format"`
	PathColumnWidth      int                   `json:"path-column-width"`
}

// Repository : A repository in the config
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

// NewConfig : Create a new config readergit diff
func NewConfig() *Config {
	return new(Config)
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
	if err != nil {
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
	_, _ = configFile.Write(data)

	_ = configFile.Close()
}

// GetConfigPath : Get path to the config file
func (config *Config) GetConfigPath() string {
	home := config.getHomeDirectory()

	return path.Join(home, defaultConfigPath)
}

// AddRepository : Adds a new repository to the config
func (config *Config) AddRepository(path, imagePath string) {
	repo := Repository{Path: path, ImagePath: imagePath}
	config.Repositories = append(config.Repositories, repo)
}

func (config *Config) GetExternalApplicationByName(name string) *ExternalApplication {
	for i := range config.ExternalApplications {
		ext := config.ExternalApplications[i]
		if ext.Name == name {
			return &ext
		}
	}

	return nil
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
