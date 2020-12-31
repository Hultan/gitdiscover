package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

const (
	exitNormal = 0
	exitConfigError = 1
	exitArgumentError = 2
)

const (
	applicationVersion = "2.0"
)

type GitStatus struct {
	Status string
	Date   *time.Time
}

var config *Config

func main() {
	// Check command line arguments
	handled := checkArguments()
	if handled {
		os.Exit(exitNormal)
	}

	// Load config
	if tryLoadConfig() == false {
		os.Exit(exitConfigError) // Failed to load config file
	}

	// Get the git statuses of the paths in the config
	var gitStatuses []GitStatus
	for _, basePath := range config.Paths {
		gitPath := path.Join(basePath, ".git")
		status := GitStatus{}

		if _, err := os.Stat(gitPath); os.IsNotExist(err) {
			status.Date = nil
			status.Status = fmt.Sprintf(createErrorFormatString(config), basePath, err)
		} else {
			gs := getGitStatus(basePath)
			status.Date = getDirectoryModifiedDate(basePath)
			status.Status = fmt.Sprintf(createFormatString(config), basePath, gs)
		}

		gitStatuses = append(gitStatuses, status)
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
		fmt.Printf("%v - %v", status.Date.Format(config.DateFormat), status.Status)
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
		config = NewConfig()
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

// Create format string for successful git status
func createFormatString(config *Config) string {
	return "%-" + strconv.Itoa(config.PathColumnWidth) + "s : %s"
}

// Create format string for failed git statuses
func createErrorFormatString(config *Config) string {
	return "%-" + strconv.Itoa(config.PathColumnWidth) + "s : Not a git directory! err=%v"
}

// Get the git status
func getGitStatus(path string) string {
	cmd := exec.Command("/home/per/bin/gitprompt-go")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return "failed to check git status"
	}
	return string(out)
}

// Get the last modified date of any file in directory...
func getDirectoryModifiedDate(directory string) *time.Time {
	var files []string

	e := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		// On errors, move on...
		if err!=nil {
			return nil
		}
		// Skip .git and .idea directories
		if info.IsDir() && (info.Name() == ".git" || info.Name() == ".idea") {
			return filepath.SkipDir
		}
		// Add all files to slice...
		if !info.IsDir() {
			files = append(files, info.Name())
		}
		return nil
	})
	if e != nil {
		log.Fatal(e)
	}

	var date *time.Time
	for _, fileName := range files {
		modDate := getModifiedDate(path.Join(directory,fileName))
		if date == nil || (modDate!=nil && modDate.After(*date)) {
			date = modDate
		}
	}
	return date
}

// Get the modified date of a file
func getModifiedDate(path string) *time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}
	date := info.ModTime()
	return &date
}

func tryLoadConfig() bool {
	// Load config
	config = NewConfig()
	if config.ConfigExists() {
		err := config.Load()
		if err != nil {
			fmt.Println("Failed to load config file (~/.config/softteam/gitdiscovery/config.json).")
			fmt.Println(err)
			return false
		}
	} else {
		fmt.Println("Missing config file (~/.config/softteam/gitdiscovery/config.json)!")
		fmt.Print("Creating a new empty one...")
		err := config.CreateEmptyConfig()
		if err!=nil {
			fmt.Println("failed!")
			fmt.Println(err)
			return false
		}
		fmt.Println("done!")
	}

	return true
}