package util

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplateFile(t *testing.T) {
	templating := Templating{
		RuleID:       "Test Rule ID",
		RuleTitle:    "Test Rule Title",
		RuleSeverity: "low",
		Replace:      strings.ReplaceAll,
	}

	currentDirectory, err := os.Getwd()
	assert.Nil(t, err)

	testTable := []struct {
		template       string
		fileName       string
		expectedResult string
	}{
		{
			template: "templates/main.tpl.rego",
			fileName: "test.rego",
			expectedResult: `package rules

deny[msg] {
	resource := input.resource.test[name]
	resource.todo
	msg := {
		# Mandatory fields
		"publicId": "Test Rule ID",
		"title": "Test Rule Title",
		"severity": "low",
		"msg": sprintf("input.resource.test[%s].todo", [name]),
		# Optional fields
		"issue": "",
		"impact": "",
		"remediation": "",
		"references": [],
	}
}
`,
		},
		{
			template: "templates/main_test.tpl.rego",
			fileName: "test.rego",
			expectedResult: `package rules

import data.lib
import data.lib.testing

test_Test Rule ID {
	# array containing test cases where the rule is allowed
	allowed_test_cases := [{
		"want_msgs": [],
		"fixture": "allowed.json",
	}]

	# array containing cases where the rule is denied
	denied_test_cases := [{
		"want_msgs": ["input.resource.test[denied].todo"], # verifies that the correct msg is returned by the denied rule
		"fixture": "denied2.tf",
	}, {
		"want_msgs": ["input.resource.test[denied].todo"], # verifies that the correct msg is returned by the denied rule
		"fixture": "denied1.yaml",
	}]

	test_cases := array.concat(allowed_test_cases, denied_test_cases)
	testing.evaluate_test_cases("Test Rule ID", "./rules/Test Rule ID/fixtures", test_cases)
}
`,
		},
		{
			template: "templates/lib/main.tpl.rego",
			fileName: "test.rego",
			expectedResult: `# File was generated automatically by the snyk-iac-rules tool
# It contains the utility functions for writing Snyk custom rules - modify at your own risk!
package lib

has_field(obj, field) {
	_ := obj[field]
}

normalize_to_array(resource) = out_array {
	is_array(resource)
	out_array = resource
} else = out_array {
	out_array = [resource]
}

merge_objects(a, b) = c {
	keys := {k | some k; _ = a[k]} | {k | some k; _ = b[k]}
	c := {k: v | k := keys[_]; v := pick(k, b, a)}
}

pick(k, obj1, _) = obj1[k]

pick(k, obj1, obj2) = obj2[k] {
	not has_field(obj1, k)
}

normalize_to_array(resource) = out_array {
	is_array(resource)
	out_array = resource
} else = out_array {
	out_array = [resource]
}
`,
		},
		{
			template: "templates/lib/testing/main.tpl.rego",
			fileName: "test.rego",
			expectedResult: `# File was generated automatically by the snyk-iac-rules tool
# It contains the testing framework for Snyk custom rules - modify at your own risk!
package lib.testing

import data.lib

assert_response_set(result_set, test_case) {
	total_violations := {res |
		result := result_set[index]
		result.publicId == test_case.publicId
		trace(sprintf("[%s][%d] Issue msg : %s", [test_case.publicId, test_case.index, result.msg]))
		res := index
	}

	trace(sprintf("[%s][%s] Number of issues identified: want %d, got %d", [test_case.publicId, test_case.fixture, count(test_case.want_msgs), count(total_violations)]))
	count(total_violations) == count(test_case.want_msgs)

	violation_match := {res |
		result := total_violations[index]
		result.msg == test_case.want_msgs[_]
		trace(sprintf("[%s][%d] Violation msg : %s", [test_case.publicId, test_case.index, result.msg]))
		res := index
	}

	trace(sprintf("[%s][%s] Number of issues with correct ` + "`msg`" + ` value: want %d, got %d", [test_case.publicId, test_case.fixture, count(test_case.want_msgs), count(violation_match)]))
	count(violation_match) == count(test_case.want_msgs)
	trace(sprintf("[%s] Fixture %d passed", [test_case.publicId, test_case.index]))
} else = false {
	true
}

parse_fixture_file(fixture_file) = fixture {
	endswith(fixture_file, "yaml")
	fixture := lib.normalize_to_array(yaml.unmarshal_file(fixture_file))
} else = fixture {
	endswith(fixture_file, "yml")
	fixture := lib.normalize_to_array(yaml.unmarshal_file(fixture_file))
} else = fixture {
	endswith(fixture_file, "tf")
	fixture := lib.normalize_to_array(hcl2.unmarshal_file(fixture_file))
} else = fixture {
	endswith(fixture_file, "json")
	fixture := lib.normalize_to_array(yaml.unmarshal_file(fixture_file))
}

get_result_set(fixture) = result_set {
	result_set := data.rules.deny with input as fixture
}

evaluate_test_cases(publicId, fixture_directory, test_cases) {
	passed_tests := {res |
		tc := lib.merge_objects(test_cases[index], {"publicId": publicId, "index": index})
		fixtures := parse_fixture_file(sprintf("%s/%s", [fixture_directory, tc.fixture]))
		result_set := get_result_set(fixtures[doc_id])
		assert_response_set(result_set, tc)
		res := index
	}

	trace(sprintf("[%s] Number of test cases passed: want %d, got %d", [publicId, count(test_cases), count(passed_tests)]))
	count(passed_tests) == count(test_cases)
} else = false {
	true
}
`,
		},
	}
	for _, test := range testTable {
		t.Run(test.template, func(t *testing.T) {
			err := TemplateFile(currentDirectory, test.fileName, test.template, templating)
			assert.Nil(t, err)

			data, err := os.ReadFile(path.Join(currentDirectory, test.fileName))
			assert.Nil(t, err)
			assert.Equal(t, test.expectedResult, string(data))

			err = os.Remove(path.Join(currentDirectory, test.fileName))
			assert.Nil(t, err)
		})
	}
}
