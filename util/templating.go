package util

import (
	"embed"
	"fmt"
	"os"
	"text/template"
)

//go:embed templates/*
//go:embed templates/lib/*
//go:embed templates/lib/testing/*
var templateFs embed.FS

type Templating struct {
	RuleName string
	Replace  func(string, string, string) string
}

func TemplateFile(workingDirectory string, fileName string, template string, templating Templating) error {
	filePath, err := CreateFile(workingDirectory, fileName)
	if err != nil {
		return err
	}
	err = generateTemplate(filePath, template, templating)
	if err != nil {
		return err
	}
	fmt.Printf("Templated file %s\n", filePath)
	return nil
}

func generateTemplate(fileName string, templateFile string, templating Templating) error {
	tpl, err := template.ParseFS(templateFs, templateFile)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	return tpl.Execute(file, templating)
}
