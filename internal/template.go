package internal

import (
	"errors"
	"fmt"
	"os"
	"path"
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

	err = templateFile(ruleFixtureDir, "allowed.json", "templates/fixtures/allowed.json", templating)
	if err != nil {
		return err
	}
	fmt.Printf("[/] Template rules/%s/fixtures/allowed.json file\n", templating.RuleID)

	err = templateFile(ruleFixtureDir, "denied1.yaml", "templates/fixtures/denied1.yaml", templating)
	if err != nil {
		return err
	}
	fmt.Printf("[/] Template rules/%s/fixtures/denied1.yaml file\n", templating.RuleID)

	err = templateFile(ruleFixtureDir, "denied2.tf", "templates/fixtures/denied2.tf", templating)
	if err != nil {
		return err
	}
	fmt.Printf("[/] Template rules/%s/fixtures/denied2.tf file\n", templating.RuleID)

	err = templateFile(ruleFixtureDir, "denied.json.tfplan", "templates/fixtures/denied.json.tfplan", templating)
	if err != nil {
		return err
	}
	fmt.Printf("[/] Template rules/%s/fixtures/denied.json.tfplan file\n", templating.RuleID)

	return nil
}

var tfPlanExists = func(tfPlan string) bool {
	_, err := os.Stat(tfPlan)
	return !errors.Is(err, os.ErrNotExist)
}

func templateLib(workingDirectory string, templating util.Templating) error {
	var libDir string
	var testingDir string
	var err error

	libDir, err = createDirectory(workingDirectory, "lib", true)
	if err != nil {
		if strings.Contains(err.Error(), "Directory already exists at") {
			// We have added a new testing helper so we should check for that first
			testingDir := path.Join(path.Join(workingDirectory, "lib"), "testing")
			if !tfPlanExists(path.Join(testingDir, "tfplan.rego")) {
				err = templateFile(testingDir, "tfplan.rego", "templates/lib/testing/tfplan.tpl.rego", templating)
				if err == nil {
					fmt.Println("[/] Template lib/testing/tfplan.rego file")
				}
			}
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
		return err
	}
	fmt.Println("[/] Template lib/testing directory")

	err = templateFile(testingDir, "main.rego", "templates/lib/testing/main.tpl.rego", templating)
	if err != nil {
		return err
	}
	fmt.Println("[/] Template lib/testing/main.rego file")

	err = templateFile(testingDir, "tfplan.rego", "templates/lib/testing/tfplan.tpl.rego", templating)
	if err != nil {
		return err
	}
	fmt.Println("[/] Template lib/testing/tfplan.rego file")
	return nil
}
