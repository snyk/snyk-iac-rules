package internal

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/snyk/snyk-iac-rules/util"
)

type TemplateCommandParams struct {
	RuleID       string
	RuleTitle    string
	RuleSeverity util.EnumFlag
	RuleFormat   util.EnumFlag
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
	err := templateRule(workingDirectory, templating, params.RuleFormat.String())
	if err != nil {
		return err
	}

	return templateLib(workingDirectory, templating)
}

func getTemplateByFormat(format string) (string, string, error) {
	switch format {
	case util.JSON:
		return "main_test_json", ".json", nil
	case util.YAML:
		return "main_test_yaml", ".yaml", nil
	case util.HCL2:
		return "main_test_hcl2", ".tf", nil
	case util.TERRAFORM_PLAN:
		return "main_test_tfplan", ".json.tfplan", nil
	default:
		// should never get to here
		return "", "", fmt.Errorf("Provided format not supported: %s", format)
	}
}

func templateRule(workingDirectory string, templating util.Templating, format string) error {
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

	testTemplate, fixtureFileExtension, err := getTemplateByFormat(format)
	if err != nil {
		return err
	}

	err = templateFile(ruleDir, "main_test.rego", "templates/"+testTemplate+".tpl.rego", templating)
	if err != nil {
		return err
	}
	fmt.Printf("[/] Template rules/%s/main_test.rego file\n", templating.RuleID)

	ruleFixtureDir, err = createDirectory(ruleDir, "fixtures", true)
	if err != nil {
		return err
	}
	fmt.Printf("[/] Template rules/%s/fixtures directory\n", templating.RuleID)

	err = templateFile(ruleFixtureDir, "denied"+fixtureFileExtension, "templates/fixtures/denied"+fixtureFileExtension, templating)
	if err != nil {
		return err
	}
	fmt.Printf("[/] Template rules/%s/fixtures/denied%s file\n", templating.RuleID, fixtureFileExtension)

	err = templateFile(ruleFixtureDir, "allowed"+fixtureFileExtension, "templates/fixtures/allowed"+fixtureFileExtension, templating)
	if err != nil {
		return err
	}
	fmt.Printf("[/] Template rules/%s/fixtures/allowed%s file\n", templating.RuleID, fixtureFileExtension)

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
