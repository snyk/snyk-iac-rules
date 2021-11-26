package util

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
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

func ValidateFilePath(path string) (fs.FileInfo, error) {
	invalidFilePath := errors.New("Failed to read from the provided path")

	file, err := os.Open(path)
	if err != nil {
		return nil, invalidFilePath
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, invalidFilePath
	}
	return fileInfo, nil
}

func CheckIfRunningInRootDirectory(paths []string) {
	for _, providedPath := range paths {
		// the user can provide the path to the rules folder
		if strings.HasSuffix(providedPath, "/rules") {
			return
		}
		// otherwise, check that the path contains the rules folder
		_, err := os.Stat(computePath(providedPath, "rules/"))
		if err == nil || !os.IsNotExist(err) {
			return
		}
	}
	// TODO: what if not templated
	fmt.Println("WARNING: The command must point at a folder that contains the package for the rules.\n	 If the rules were generated using the template command, make sure you have the\n	 /rules and /lib folder in your current running directory or provide an optional path argument pointing to that location.")
}
