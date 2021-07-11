package main

import (
	"errors"
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/gitdiscover/internal/config"
	"github.com/hultan/gitdiscover/internal/gitdiscover"
	"github.com/hultan/gitdiscover/internal/gui"
	"os"
	"sort"
)

func main() {
	gui := isGuiRequested()
	if gui {
		fmt.Printf("Gitdiscover : Starting GUI!\n")
		showGUI()
		return
	}

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
	fmt.Println("Git Repository Statuses")
	fmt.Println("_______________________")
	for _, status := range gitStatuses {
		var text = ""
		if status.Date == nil {
			text = fmt.Sprintf("%v - %v - %v", "2006-01-02, kl. 15:04", status.Path, status.Status)
		} else {
			text = fmt.Sprintf("%s - %s - %s", status.Date.Format(config.DateFormat), status.Path, status.Status)
		}
		fmt.Println(text)
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

func isGuiRequested() bool {
	if len(os.Args) == 1 {
		return false
	}

	if os.Args[1] == "-gui" {
		return true
	}

	return false
}

func showGUI() {
	// Create a new application
	application, err := gtk.ApplicationNew(ApplicationId, ApplicationFlags)
	if err != nil {
		panic(err)
	}

	mainForm := gui.NewMainWindow()
	// Hook up the activate event handler
	_ = application.Connect("activate", mainForm.OpenMainWindow)

	// Start the application (and exit when it is done)
	os.Exit(application.Run(nil))
}
