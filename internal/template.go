package internal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/snyk/snyk-iac-custom-rules/util"
)

type TemplateCommandParams struct {
	Rule string
}

var createDirectory = util.CreateDirectory
var templateFile = util.TemplateFile

func RunTemplate(args []string, params *TemplateCommandParams) error {
	workingDirectory := args[0]

	templating := util.Templating{
		RuleName: params.Rule,
		Replace:  strings.ReplaceAll,
	}
	err := templateRule(workingDirectory, templating)
	if err != nil {
		return err
	}

	return templateLib(workingDirectory, templating)
}

func templateRule(workingDirectory string, templating util.Templating) error {
	rulesDir, err := createDirectory(workingDirectory, "rules", false)
	if err != nil {
		return err
	}
	fmt.Printf("Templated directory %s\n", rulesDir)

	ruleDir, err := createDirectory(rulesDir, templating.RuleName, true)
	if err != nil {
		if strings.Contains(err.Error(), "Directory already exists") {
			return errors.New("Rule with the provided name already exists")
		}
		return err
	}
	fmt.Printf("Templated directory %s\n", ruleDir)

	err = templateFile(ruleDir, "main.rego", "templates/main.tpl.rego", templating)
	if err != nil {
		return err
	}

	err = templateFile(ruleDir, "main_test.rego", "templates/main_test.tpl.rego", templating)
	if err != nil {
		return err
	}

	return nil
}

func templateLib(workingDirectory string, templating util.Templating) error {
	libDir, err := createDirectory(workingDirectory, "lib", true)
	if err != nil {
		if strings.Contains(err.Error(), "Directory already exists at") {
			return nil
		}
		return err
	}
	fmt.Printf("Templated directory %s\n", libDir)

	err = templateFile(libDir, "main.rego", "templates/lib/main.tpl.rego", templating)
	if err != nil {
		return err
	}

	testingDir, err := createDirectory(libDir, "testing", true)
	if err != nil {
		if strings.Contains(err.Error(), "Directory already exists at") {
			return nil
		}
		return err
	}
	fmt.Printf("Templated directory %s\n", testingDir)

	err = templateFile(testingDir, "main.rego", "templates/lib/testing/main.tpl.rego", templating)
	if err != nil {
		return err
	}
	return nil
}
