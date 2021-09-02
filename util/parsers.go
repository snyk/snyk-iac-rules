package util

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
	hcl2jsonConverter "github.com/tmccombs/hcl2json/convert"
)

func ParseHCL2(p []byte, v interface{}) error {
	jsonBytes, err := hcl2jsonConverter.Bytes(p, "", hcl2jsonConverter.Options{})
	if err != nil {
		return fmt.Errorf("convert to bytes: %w", err)
	}

	if err := json.Unmarshal(jsonBytes, v); err != nil {
		return fmt.Errorf("unmarshal hcl2: %w", err)
	}
	return nil
}

// ParseYAML unmarshals YAML files and return parsed file content.
func ParseYAML(p []byte, v interface{}) error {
	subDocuments := separateSubDocuments(p)
	if len(subDocuments) > 1 {
		if err := unmarshalMultipleDocuments(subDocuments, v); err != nil {
			return fmt.Errorf("unmarshal multiple documents: %w", err)
		}

		return nil
	}

	if err := yaml.Unmarshal(p, v); err != nil {
		return fmt.Errorf("unmarshal yaml: %w", err)
	}

	return nil
}

func separateSubDocuments(data []byte) [][]byte {
	linebreak := "\n"
	if bytes.Contains(data, []byte("\r\n---\r\n")) {
		linebreak = "\r\n"
	}

	return bytes.Split(data, []byte(linebreak+"---"+linebreak))
}

func unmarshalMultipleDocuments(subDocuments [][]byte, v interface{}) error {
	var documentStore []interface{}
	for _, subDocument := range subDocuments {
		var documentObject interface{}
		if err := yaml.Unmarshal(subDocument, &documentObject); err != nil {
			return fmt.Errorf("unmarshal subdocument yaml: %w", err)
		}

		documentStore = append(documentStore, documentObject)
	}

	yamlConfigBytes, err := yaml.Marshal(documentStore)
	if err != nil {
		return fmt.Errorf("marshal yaml document: %w", err)
	}

	if err := yaml.Unmarshal(yamlConfigBytes, v); err != nil {
		return fmt.Errorf("unmarshal yaml: %w", err)
	}

	return nil
}
