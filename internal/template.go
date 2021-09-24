package internal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/snyk/snyk-iac-custom-rules/util"
)

const (
	LOW      = "low"
	MEDIUM   = "medium"
	HIGH     = "high"
	CRITICAL = "critical"
)

type TemplateCommandParams struct {
	RuleID       string
	RuleTitle    string
	RuleSeverity util.EnumFlag
}

var createDirectory = util.CreateDirectory
var templateFile = util.TemplateFile
var startProgress = util.StartProgress

func RunTemplate(args []string, params *TemplateCommandParams) error {
	workingDirectory := args[0]

	templating := util.Templating{
		RuleID:       params.RuleID,
		RuleTitle:    params.RuleTitle,
		RuleSeverity: params.RuleSeverity.String(),
		Replace:      strings.ReplaceAll,
	}
	err := templateRule(workingDirectory, templating)
	if err != nil {
		return err
	}

	return templateLib(workingDirectory, templating)
}

func templateRule(workingDirectory string, templating util.Templating) error {
	var rulesDir string
	var ruleDir string
	var ruleFixtureDir string
	var err error

	err = startProgress("Template rules directory", func() error {
		rulesDir, err = createDirectory(workingDirectory, "rules", false)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = startProgress(fmt.Sprintf("Template rules/%s directory", templating.RuleID), func() error {
		ruleDir, err = createDirectory(rulesDir, templating.RuleID, true)
		if err != nil {
			if strings.Contains(err.Error(), "Directory already exists") {
				return errors.New("Rule with the provided name already exists")
			}
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = startProgress(fmt.Sprintf("Template rules/%s/main.rego file", templating.RuleID), func() error {
		return templateFile(ruleDir, "main.rego", "templates/main.tpl.rego", templating)
	})
	if err != nil {
		return err
	}

	err = startProgress(fmt.Sprintf("Template rules/%s/main_test.rego file", templating.RuleID), func() error {
		return templateFile(ruleDir, "main_test.rego", "templates/main_test.tpl.rego", templating)
	})
	if err != nil {
		return err
	}

	err = startProgress(fmt.Sprintf("Template rules/%s/fixtures directory", templating.RuleID), func() error {
		ruleFixtureDir, err = createDirectory(ruleDir, "fixtures", true)
		return nil
	})
	if err != nil {
		return err
	}

	err = startProgress(fmt.Sprintf("Template rules/%s/fixtures/allowed.json file", templating.RuleID), func() error {
		return templateFile(ruleFixtureDir, "allowed.json", "templates/fixtures/allowed.json", templating)
	})
	if err != nil {
		return err
	}

	err = startProgress(fmt.Sprintf("Template rules/%s/fixtures/denied.json file", templating.RuleID), func() error {
		return templateFile(ruleFixtureDir, "denied.json", "templates/fixtures/denied.json", templating)
	})
	if err != nil {
		return err
	}

	return nil
}

func templateLib(workingDirectory string, templating util.Templating) error {
	var libDir string
	var testingDir string
	var err error

	err = startProgress("Template lib directory", func() error {
		libDir, err = createDirectory(workingDirectory, "lib", true)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "Directory already exists at") {
			return nil
		}
		return err
	}

	err = startProgress("Template lib/main.rego file", func() error {
		return templateFile(libDir, "main.rego", "templates/lib/main.tpl.rego", templating)
	})
	if err != nil {
		return err
	}

	err = startProgress("Template lib/testing directory", func() error {
		testingDir, err = createDirectory(libDir, "testing", true)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "Directory already exists at") {
			return nil
		}
		return err
	}

	return startProgress("Template lib/testing/main.rego file", func() error {
		return templateFile(testingDir, "main.rego", "templates/lib/testing/main.tpl.rego", templating)
	})
}
