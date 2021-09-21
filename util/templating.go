package util

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"text/template"
)

//go:embed templates/*
//go:embed templates/lib/*
//go:embed templates/lib/testing/*
var templateFs embed.FS

type Templating struct {
	RuleID       string
	RuleTitle    string
	RuleSeverity string
	Replace      func(string, string, string) string
}

var createFile = CreateFile
var openFile = func(name string, flag int, mode fs.FileMode) (io.Writer, error) {
	return os.OpenFile(name, flag, mode)
}

func TemplateFile(workingDirectory string, fileName string, template string, templating Templating) error {
	filePath, err := createFile(workingDirectory, fileName)
	if err != nil {
		return err
	}
	err = generateTemplate(filePath, template, templating)
	if err != nil {
		return err
	}
	return nil
}

func generateTemplate(filePath string, templateFile string, templating Templating) error {
	tpl, err := template.ParseFS(templateFs, templateFile)
	if err != nil {
		return err
	}

	file, err := openFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	return tpl.Execute(file, templating)
}
