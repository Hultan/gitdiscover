package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/gotk3/gotk3/gtk"
	"github.com/sirupsen/logrus"

	gitConfig "github.com/hultan/gitdiscover/internal/config"
	"github.com/hultan/gitdiscover/internal/gitdiscover"
	"github.com/hultan/gitdiscover/internal/gui"
)

var (
	logger *logrus.Logger
)

func main() {
	// Logging and config
	logger = startLogging()
	config := loadConfig()

	// Check if GUI is reuqested
	if len(os.Args) > 1 {
		logger.Info("Starting GitDiscover GUI!")
		showGUI(logger, config)
	}

	// Get repository list
	git := gitdiscover.NewGit(config, logger)

	// Sort the git status string after modified date of the .git folder
	sort.Slice(git.Repos, func(i, j int) bool {
		date1 := git.Repos[i].ModifiedDate
		date2 := git.Repos[j].ModifiedDate
		if date1 == nil || date2 == nil {
			return false
		}
		return (*date1).After(*date2)
	})

	// Print out the git statuses
	fmt.Println("Git Repository Statuses : ")
	logger.Info("Git Repository Statuses : ")
	fmt.Println("_______________________")
	for _, status := range git.Repos {
		var text = ""
		if status.ModifiedDate == nil {
			text = fmt.Sprintf("%v - %v - %v", "2006-01-02, kl. 15:04", status.Path, status.Status)
		} else {
			text = fmt.Sprintf("%s - %s - %s", status.ModifiedDate.Format(config.DateFormat), status.Path, status.Status)
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
