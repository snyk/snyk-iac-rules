package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	parsers "github.com/snyk/snyk-iac-parsers"

	"github.com/snyk/snyk-iac-rules/util"
)

var readFile = os.ReadFile
var parseYAML = parsers.ParseYAML
var parseHCL2 = parsers.ParseHCL2
var parseTerraformPlan = parsers.ParseTerraformPlan

type ParseCommandParams struct {
	Format util.EnumFlag
}

func RunParse(args []string, params *ParseCommandParams) error {
	filePath := args[0]

	content, err := readFile(filePath)
	if err != nil {
		return err
	}

	var parsedInput interface{}
	switch params.Format.String() {
	case util.YAML:
		if err := parseYAML(content, &parsedInput); err != nil {
			return err
		}
	case util.TERRAFORM_PLAN:
		if err := parseTerraformPlan(content, &parsedInput); err != nil {
			return err
		}
	default:
		if err := parseHCL2(content, &parsedInput); err != nil {
			return err
		}
	}

	jsonInput, err := json.Marshal(parsedInput)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	if err := json.Indent(&out, jsonInput, "", "\t"); err != nil {
		return err
	}

	fmt.Printf("%s\n", out.String())
	return nil
}
