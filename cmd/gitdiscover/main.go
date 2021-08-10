package main

import (
	"errors"
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	gitConfig "github.com/hultan/gitdiscover/internal/config"
	"github.com/hultan/gitdiscover/internal/gitdiscover"
	"github.com/hultan/gitdiscover/internal/gui"
	"github.com/sirupsen/logrus"
	"os"
	"sort"
)

var (
	logger *logrus.Logger
)

func main() {
	// Logging and config
	logger = startLogging()
	config := loadConfig()

	// Check if GUI is reuqested
	guiRequested := isGuiRequested()
	if guiRequested {
		logger.Info("Starting GitDiscover GUI!")
		showGUI(logger, config)
	}

	// Check command line arguments
	checkArguments()

	// Get repository list
	gitDiscover := gitdiscover.NewGit(config, logger)
	gitStatuses, err := gitDiscover.GetRepositories()
	if err != nil {
		exitProgram(exitUnknown, err)
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
		fmt.Println(text)
		logger.Info(text)
	}

	exitProgram(exitNormal, nil)
}

func exitProgram(exitCode int, err error) {
	if err!=nil {
		fmt.Println(err)
		logger.Error(err)
	}
	stopLogging(logger)
	os.Exit(exitCode)
}

//
// Logging functions
//

func stopLogging(logger *logrus.Logger) {
	logger.Info("Exit GitDiscover")
	logger = nil
}

func startLogging() *logrus.Logger {
	l := logrus.New()
	l.Level = logrus.TraceLevel
	l.Out = os.Stdout

	file, err := os.OpenFile(ApplicationLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		l.Out = file
	} else {
		l.Info("Failed to log to file, using default stderr")
	}
	l.Info("Starting GitDiscover")
	return l
}

//
// Config functions
//

func loadConfig() *gitConfig.Config {
	config := gitConfig.NewConfig()
	if config.ConfigExists() {
		// Existing config file
		err := config.Load()
		if err != nil {
			exitProgram(exitConfigError, err)
		}
	} else {
		// New config file
		err := config.CreateEmptyConfig()
		if err!=nil {
			exitProgram(exitConfigError, err)
		}
		fmt.Println("done!")
	}

	return config
}

//
// GUI functions
//

func isGuiRequested() bool {
	if len(os.Args) == 1 {
		return false
	}

	if os.Args[1] == "-gui" {
		return true
	}

	return false
}

func showGUI(logger *logrus.Logger, config *gitConfig.Config) {
	// Create a new application
	application, err := gtk.ApplicationNew(ApplicationId, ApplicationFlags)
	if err != nil {
		exitProgram(exitUnknown, err)
	}

	mainForm := gui.NewMainWindow(logger, config)
	mainForm.ApplicationLogPath = ApplicationLogPath
	// Hook up the activate event handler
	_ = application.Connect("activate", mainForm.OpenMainWindow)

	// Start the application (and exit when it is done)
	exitCode := application.Run(nil)
	exitProgram(exitCode, nil)
}

//
// Check command line arguments
//

func checkArguments() {
	if len(os.Args) > 1 {
		if os.Args[1] == "--version" {
			version := gui.ApplicationVersion
			logger.Info("Version requested : ", version)
			fmt.Printf("Gitdiscover %s\n", version)
			exitProgram(exitNormal, nil)
		} else if os.Args[1] == "--help" {
			logger.Info("Help requested.")
			fmt.Println("Usage : gitdiscover [--version] [--help]")
			exitProgram(exitNormal, nil)
		} else {
			err := errors.New(fmt.Sprintf("Invalid argument : %s", os.Args[1]))
			fmt.Println("Usage : gitdiscover [--version] [--help]")
			exitProgram(exitArgumentError, err)
		}
	}
}
