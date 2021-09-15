package util

import (
	"fmt"
	"os"
	"path"
)

func checkIfDirectoryExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

var createDirectoryIfNotExists = func(dirPath string, strict bool) error {
	exists, _ := checkIfDirectoryExists(dirPath)
	if !exists {
		if err := os.Mkdir(dirPath, 0755); err != nil {
			return err
		}
	} else if strict {
		return fmt.Errorf("Directory already exists at path %s", dirPath)
	}
	return nil
}

func CreateDirectory(workingDirectory string, name string, strict bool) (string, error) {
	dirPath := computePath(workingDirectory, name)
	err := createDirectoryIfNotExists(dirPath, strict)
	return dirPath, err
}

var createFileAndClose = func(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

var computePath = func(workingDirectory string, name string) string {
	return path.Join(workingDirectory, name)
}

func CreateFile(workingDirectory string, name string) (string, error) {
	filePath := computePath(workingDirectory, name)
	err := createFileAndClose(filePath)
	if err != nil {
		return filePath, err
	}
	return filePath, nil
}
