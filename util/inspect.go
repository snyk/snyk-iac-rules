package util

import (
	"github.com/open-policy-agent/opa/loader"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func RetrieveRules(paths []string) ([]string, error) {
	fileNames, err := findRegoFiles(paths)
	if err != nil {
		return []string{}, err
	}

	var publicIds = []string{}
	for _, fileName := range fileNames {
		publicId, err := getPublicIdFromFile(fileName)
		if err != nil {
			return []string{}, err
		}
		if publicId != "" {
			publicIds = append(publicIds, publicId)
		}
	}
	return publicIds, nil
}

func findRegoFiles(paths []string) ([]string, error) {
	fileNames := []string{}

	result, err := loader.AllRegos(paths)
	if err != nil {
		return fileNames, err
	}

	for _, module := range result.Modules {
		fileName := filepath.Clean(module.Name)
		if !strings.Contains(fileName, "_test.rego") {
			fileNames = append(fileNames, fileName)
		}
	}

	return fileNames, nil
}

func getPublicIdFromFile(fileName string) (string, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	publicId := extractPublicIdFromRego(string(data))
	return publicId, nil
}

func extractPublicIdFromRego(rego string) string {
	re := regexp.MustCompile("\"publicId\"\\s*:\\s*\"(.*?)\"")
	match := re.FindStringSubmatch(rego)
	if len(match) > 0 {
		return match[1]
	}
	return ""
}
