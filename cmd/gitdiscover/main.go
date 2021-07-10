package main

import (
	"errors"
	"fmt"
	"github.com/hultan/gitdiscover/internal/config"
	"github.com/hultan/gitdiscover/internal/gitdiscover"
	"os"
	"sort"
)

func main() {
	// Check command line arguments
	handled := checkArguments()
	if handled {
		os.Exit(exitNormal)
	}

	// Load config
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	gitDiscover := gitdiscover.GitNew(config)
	gitStatuses, err := gitDiscover.GetRepositories()
	if err != nil {
		panic(err)
	}

	// Sort the git status string after modified date of the .git folder
	sort.Slice(gitStatuses, func(i, j int) bool {
		date1 := gitStatuses[i].Date
		date2 := gitStatuses[j].Date
		if date1 == nil || date2 == nil {
			return false
		}
		return (*date1).After(*date2)
	})

	// Print out the git statuses
	for _, status := range gitStatuses {
		if status.Date == nil {
			fmt.Printf("%v - %v\n", "2006-01-02, kl. 15:04", status.Status)
		} else {
			fmt.Printf("%v - %v", status.Date.Format(config.DateFormat), status.Status)
		}
	}

	// Exit
	os.Exit(exitNormal)
}

func checkArguments() bool {
	if len(os.Args) == 1 {
		return false
	}

	if os.Args[1] == "--version" {
		fmt.Printf("Gitdiscover %s\n", applicationVersion)
		return true
	} else if os.Args[1] == "--help" {
		fmt.Println("Usage : gitdiscover [--version] [--help] [add-path path]")
		return true
	} else if os.Args[1] == "add-path" && len(os.Args) == 3 {
		config := config.NewConfig()
		err:=config.Load()
		if err!=nil {
			fmt.Println(err)
			os.Exit(exitConfigError)
		}
		config.Paths = append(config.Paths, os.Args[2])
		config.Save()
		fmt.Printf("The path '%s' has been added to the gitdiscover config!", os.Args[2])
		return true
	}

	fmt.Println("Invalid argument!")
	os.Exit(exitArgumentError)
	return false
}

func loadConfig() (*config.Config, error) {
	// Load config
	config := config.NewConfig()
	if config.ConfigExists() {
		err := config.Load()
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to load config file (%s).\n", config.GetConfigPath()))
		}
	} else {
		err := config.CreateEmptyConfig()
		if err!=nil {
			return nil, errors.New(fmt.Sprintf("Missing config file (%s).\n", config.GetConfigPath()))
		}
		fmt.Println("done!")
	}

	return config, nil
}