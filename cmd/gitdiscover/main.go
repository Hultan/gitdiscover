package main

import (
	"errors"
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/gitdiscover/internal/config"
	"github.com/hultan/gitdiscover/internal/gitdiscover"
	"github.com/hultan/gitdiscover/internal/gui"
	"github.com/sirupsen/logrus"
	"os"
	"sort"
)

func main() {
	logger := startLogging()

	// Check command line arguments
	guiRequested := isGuiRequested()
	if guiRequested {
		logger.Info("Starting GitDiscover GUI!")
		showGUI(logger)
		return
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "--version" {
			version := gui.GetVersion()
			logger.Info("Version requested : ", version)
			fmt.Printf("Gitdiscover %s\n", version)
			os.Exit(exitNormal)
		} else if os.Args[1] == "--help" {
			logger.Info("Help requested.")
			fmt.Println("Usage : gitdiscover [--version] [--help]")
			os.Exit(exitNormal)
		} else {
			logger.Info("Invalid argument : ", os.Args[1])
			fmt.Println("Invalid argument : ", os.Args[1])
			fmt.Println("Usage : gitdiscover [--version] [--help]")
			os.Exit(exitArgumentError)
		}
	}

	// Load config
	config, err := loadConfig()
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	gitDiscover := gitdiscover.NewGit(config, logger)
	gitStatuses, err := gitDiscover.GetRepositories()
	if err != nil {
		logger.Error(err)
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
	fmt.Println("Git Repository Statuses : ")
	logger.Info("Git Repository Statuses : ")
	fmt.Println("_______________________")
	for _, status := range gitStatuses {
		var text = ""
		if status.Date == nil {
			text = fmt.Sprintf("%v - %v - %v", "2006-01-02, kl. 15:04", status.Path, status.Status)
		} else {
			text = fmt.Sprintf("%s - %s - %s", status.Date.Format(config.DateFormat), status.Path, status.Status)
		}
		logger.Info(text)
		fmt.Println(text)
	}

	stopLogging(logger)

	// Exit
	os.Exit(exitNormal)
}

func stopLogging(logger *logrus.Logger) {
	logger.Info("Exit GitDiscover")

	logger = nil
}

func startLogging() *logrus.Logger {
	logger := logrus.New()
	logger.Level = logrus.TraceLevel
	logger.Out = os.Stdout

	file, err := os.OpenFile("logrus.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		logger.Out = file
	} else {
		logger.Info("Failed to log to file, using default stderr")
	}
	logger.Info("Starting GitDiscover")
	return logger
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

func showGUI(logger *logrus.Logger) {
	// Create a new application
	application, err := gtk.ApplicationNew(ApplicationId, ApplicationFlags)
	if err != nil {
		panic(err)
	}

	mainForm := gui.NewMainWindow(logger)
	// Hook up the activate event handler
	_ = application.Connect("activate", mainForm.OpenMainWindow)

	// Start the application (and exit when it is done)
	exitCode := application.Run(nil)
	stopLogging(logger)
	os.Exit(exitCode)
}
