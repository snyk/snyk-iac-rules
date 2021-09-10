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
	rulesDir, err := util.CreateDirectory(workingDirectory, "rules", false)
	if err != nil {
		return err
	}
	fmt.Printf("Templated directory %s\n", rulesDir)

	ruleDir, err := util.CreateDirectory(rulesDir, templating.RuleName, true)
	if err != nil {
		if strings.Contains(err.Error(), "Directory already exists") {
			return errors.New("Rule with the provided name already exists")
		}
		return err
	}
	fmt.Printf("Templated directory %s\n", ruleDir)

	err = util.TemplateFile(ruleDir, "main.rego", "templates/main.tpl.rego", templating)
	if err != nil {
		return err
	}

	err = util.TemplateFile(ruleDir, "main_test.rego", "templates/main_test.tpl.rego", templating)
	if err != nil {
		return err
	}

	return nil
}

func templateLib(workingDirectory string, templating util.Templating) error {
	libDir, err := util.CreateDirectory(workingDirectory, "lib", true)
	if err != nil {
		if strings.Contains(err.Error(), "Directory already exists at") {
			return nil
		}
		return err
	}
	fmt.Printf("Templated directory %s\n", libDir)

	err = util.TemplateFile(libDir, "main.rego", "templates/lib/main.tpl.rego", templating)
	if err != nil {
		return err
	}

	testingDir, err := util.CreateDirectory(libDir, "testing", true)
	if err != nil {
		if strings.Contains(err.Error(), "Directory already exists at") {
			return nil
		}
		return err
	}
	fmt.Printf("Templated directory %s\n", testingDir)

	err = util.TemplateFile(testingDir, "main.rego", "templates/lib/testing/main.tpl.rego", templating)
	if err != nil {
		return err
	}
	return nil
}
