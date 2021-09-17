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
		RuleName: "Test Rule",
		Replace:  strings.ReplaceAll,
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
	input.spec.template.todo
	msg := {
		"publicId": "Test Rule",
		"title": "<TODO>",
		"severity": "<TODO>",
		"issue": "",
		"impact": "",
		"remediation": "",
		"msg": "spec.template.todo",
		"references": [],
	}
}`,
		},
		{
			template: "templates/main_test.tpl.rego",
			fileName: "test.rego",
			expectedResult: `package rules

import data.lib
import data.lib.testing

test_Test Rule {
	test_cases := [{
		"want_msgs": [],
		"fixture": {
			"spec": {
				"template": {
					"todo": false
				}
			}
		},
	}, {
		"want_msgs": ["spec.template.todo"],
		"fixture": {
			"spec": {
				"template": {
					"todo": true
				}
			}
		},
	}]

	testing.evaluate_test_cases("Test Rule", test_cases)
}`,
		},
		{
			template: "templates/lib/main.tpl.rego",
			fileName: "test.rego",
			expectedResult: `package lib

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
			expectedResult: `package lib.testing

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
	fixture := lib.normalize_to_array(yaml.unmarshal_file(fixture_file))
}

get_result_set(fixture) = result_set {
	result_set := data.rules.deny with input as fixture
}

evaluate_test_cases(publicId, test_cases) {
	passed_tests := {res |
		tc := lib.merge_objects(test_cases[index], {"publicId": publicId, "index": index})
		result_set := get_result_set(tc.fixture)
		assert_response_set(result_set, tc)
		res := index
	}

	trace(sprintf("[%s] Number of test cases passed: want %d, got %d", [publicId, count(test_cases), count(passed_tests)]))
	count(passed_tests) == count(test_cases)
} else = false {
	true
}`,
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