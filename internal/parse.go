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

type ParseCommandParams struct {
	Format util.EnumFlag
}

func RunParse(args []string, params *ParseCommandParams) error {
	filePath := args[0]

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var parsedInput interface{}
	switch params.Format.String() {
	case YAML:
		if err := util.ParseYAML(content, &parsedInput); err != nil {
			return err
		}
	default:
		// HCL2 is the only other option here
		if err := util.ParseHCL2(content, &parsedInput); err != nil {
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
