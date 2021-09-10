package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/snyk/snyk-iac-custom-rules/util"
)

const (
	HCL2 = "hcl2"
	YAML = "yaml"
)

var readFile = ioutil.ReadFile
var parseYAML = util.ParseYAML
var parseHCL2 = util.ParseHCL2

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
	case YAML:
		if err := parseYAML(content, &parsedInput); err != nil {
			return err
		}
	default:
		// HCL2 is the only other option here
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
