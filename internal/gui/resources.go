package gui

import (
	"errors"
	"os"
	"path"
	"path/filepath"
)

// Resources : Handles SoftTeam resources
type Resources struct {
}

func ResourcesNew() *Resources {
	return new(Resources)
}

// GetExecutablePath : Returns the path of the executable
func (r *Resources) GetExecutablePath() string {
	ex, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(ex)
}

// GetResourcesPath : Returns the resources path
func (r *Resources) GetResourcesPath() string {
	executablePath:=r.GetExecutablePath()

	var pathsToCheck []string
	pathsToCheck = append(pathsToCheck,path.Join(executablePath, "assets"))
	pathsToCheck = append(pathsToCheck,path.Join(executablePath, "../assets"))

	dir, err := r.checkPathsExists(pathsToCheck)
	if err!=nil {
		return executablePath
	}
	return dir
}

func (r *Resources) checkPathsExists(pathsToCheck []string) (string, error) {
	for _, pathToCheck := range pathsToCheck {
		if _, err := os.Stat(pathToCheck); os.IsNotExist(err) == false {
			return pathToCheck, nil
		}
	}
	return "", errors.New("paths do not exist")
}

// GetResourcePath : Gets the path for a single resource file
func (r *Resources) GetResourcePath(fileName string) string {
	resourcesPath:=r.GetResourcesPath()
	resourcePath:=path.Join(resourcesPath, fileName)
	return resourcePath
}

func getResourcePath(fileName string) (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := path.Dir(exePath)

	gladePath := path.Join(exeDir, fileName)
	if fileExists(gladePath) {
		return gladePath, nil
	}
	gladePath = path.Join(exeDir, "assets", fileName)
	if fileExists(gladePath) {
		return gladePath, nil
	}
	gladePath = path.Join(exeDir, "../assets", fileName)
	if fileExists(gladePath) {
		return gladePath, nil
	}
	return gladePath, nil
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
