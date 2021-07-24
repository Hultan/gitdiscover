package gui

import (
	"errors"
	"os"
	"path"
)

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

	return "", errors.New("failed to find resource path")
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
