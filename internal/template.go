package internal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/snyk/snyk-iac-rules/util"
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

	rulesDir, err = createDirectory(workingDirectory, "rules", false)
	if err != nil {
		return err
	}
	fmt.Println("[/] Template rules directory")

	ruleDir, err = createDirectory(rulesDir, templating.RuleID, true)
	if err != nil {
		if strings.Contains(err.Error(), "Directory already exists") {
			return errors.New("Rule with the provided name already exists")
		}
		return err
	}
	fmt.Printf("[/] Template rules/%s directory\n", templating.RuleID)

	err = templateFile(ruleDir, "main.rego", "templates/main.tpl.rego", templating)
	if err != nil {
		return err
	}
	fmt.Printf("[/] Template rules/%s/main.rego file\n", templating.RuleID)

	err = templateFile(ruleDir, "main_test.rego", "templates/main_test.tpl.rego", templating)
	if err != nil {
		return err
	}
	fmt.Printf("[/] Template rules/%s/main_test.rego file\n", templating.RuleID)

	ruleFixtureDir, err = createDirectory(ruleDir, "fixtures", true)
	if err != nil {
		return err
	}
	fmt.Printf("[/] Template rules/%s/fixtures directory\n", templating.RuleID)

	err = templateFile(ruleFixtureDir, "allowed.tf", "templates/fixtures/allowed.tf", templating)
	if err != nil {
		return err
	}
	fmt.Printf("[/] Template rules/%s/fixtures/allowed.tf file\n", templating.RuleID)

	err = templateFile(ruleFixtureDir, "denied.tf", "templates/fixtures/denied.tf", templating)
	if err != nil {
		return err
	}
	fmt.Printf("[/] Template rules/%s/fixtures/denied.tf file\n", templating.RuleID)
	return nil
}

func templateLib(workingDirectory string, templating util.Templating) error {
	var libDir string
	var testingDir string
	var err error

	libDir, err = createDirectory(workingDirectory, "lib", true)
	if err != nil {
		if strings.Contains(err.Error(), "Directory already exists at") {
			return nil
		}
		return err
	}
	fmt.Println("[/] Template lib directory")

	err = templateFile(libDir, "main.rego", "templates/lib/main.tpl.rego", templating)
	if err != nil {
		return err
	}
	fmt.Println("[/] Template lib/main.rego file")

	testingDir, err = createDirectory(libDir, "testing", true)
	if err != nil {
		if strings.Contains(err.Error(), "Directory already exists at") {
			return nil
		}
		return err
	}
	fmt.Println("[/] Template lib/testing directory")

	err = templateFile(testingDir, "main.rego", "templates/lib/testing/main.tpl.rego", templating)
	if err != nil {
		return err
	}
	fmt.Println("[/] Template lib/testing/main.rego file")
	return nil
}
