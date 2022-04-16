package main

import (
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gtk"
	"github.com/sirupsen/logrus"

	gitConfig "github.com/hultan/gitdiscover/internal/config"
	"github.com/hultan/gitdiscover/internal/gitdiscover-gui"
)

var (
	logger *logrus.Logger
)

func main() {
	// Logging and config
	logger = startLogging()
	config := loadConfig()

	logger.Info("Starting GitDiscover GUI!")
	showGUI(logger, config)
}

func exitProgram(exitCode int, err error) {
	if err != nil {
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

	file, err := os.OpenFile(applicationLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
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
	// Existing config file
	err := config.Load()
	if err != nil {
		exitProgram(exitConfigError, err)
	}
	return config
}

//
// GUI functions
//

func showGUI(logger *logrus.Logger, config *gitConfig.Config) {
	// Create a new application
	application, err := gtk.ApplicationNew(applicationId, applicationFlags)
	if err != nil {
		exitProgram(exitUnknown, err)
	}

	mainForm := gitdiscover_gui.NewMainWindow(logger, config)
	mainForm.ApplicationLogPath = applicationLogPath
	// Hook up the activate event handler
	_ = application.Connect("activate", mainForm.OpenMainWindow)

	// Start the application (and exit when it is done)
	exitCode := application.Run(nil)
	exitProgram(exitCode, nil)
}
