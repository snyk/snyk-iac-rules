package util

import (
	"fmt"
	"os"
	"path"
)

func CreateDirectory(workingDirectory string, name string, strict bool) (string, error) {
	dirName := path.Join(workingDirectory, name)
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		if err := os.Mkdir(dirName, 0755); err != nil {
			return "", err
		}
	} else if strict {
		return "", fmt.Errorf("Directory already exists at path %s", dirName)
	}
	return dirName, nil
}

func CreateFile(workingDirectory string, name string) (string, error) {
	fileName := path.Join(workingDirectory, name)
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()
	return fileName, nil
}
