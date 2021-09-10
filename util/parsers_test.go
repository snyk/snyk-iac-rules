package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This test is mostly attributed to https://github.com/tmccombs/hcl2json
func TestHCL2Parser(t *testing.T) {
	testTable := []struct {
		name           string
		controlConfigs []byte
		expectedResult interface{}
	}{
		{
			name: "simple-resources",
			controlConfigs: []byte(`
resource "aws_elastic_beanstalk_environment" "example" {
	name        = "test_environment"
	application = "testing"
	setting {
		namespace = "aws:autoscaling:asg"
		name      = "MinSize"
		value     = "1"
	}
	dynamic "setting" {
		for_each = data.consul_key_prefix.environment.var
		content {
		heredoc = <<-EOF
		This is a heredoc template.
		It references ${local.other.3}
		EOF
		simple = "${4 - 2}"
		cond = test3 > 2 ? 1: 0
		heredoc2 = <<EOF
			Another heredoc, that
			doesn't remove indentation
			${local.other.3}
			%{if true ? false : true}"gotcha"\n%{else}4%{endif}
		EOF
		loop = "This has a for loop: %{for x in local.arr}x,%{endfor}"
		namespace = "aws:elasticbeanstalk:application:environment"
		name      = setting.key
		value     = setting.value
		}
	}
	}`),
			expectedResult: map[string]interface{}(map[string]interface{}{"resource": map[string]interface{}{"aws_elastic_beanstalk_environment": map[string]interface{}{"example": map[string]interface{}{"application": "testing", "dynamic": map[string]interface{}{"setting": map[string]interface{}{"content": map[string]interface{}{"cond": "${test3 > 2 ? 1: 0}", "heredoc": "This is a heredoc template.\nIt references ${local.other.3}\n", "heredoc2": "\t\t\tAnother heredoc, that\n\t\t\tdoesn't remove indentation\n\t\t\t${local.other.3}\n\t\t\t%{if true ? false : true}\"gotcha\"\\n%{else}4%{endif}\n", "loop": "This has a for loop: %{for x in local.arr}x,%{endfor}", "name": "${setting.key}", "namespace": "aws:elasticbeanstalk:application:environment", "simple": "${4 - 2}", "value": "${setting.value}"}, "for_each": "${data.consul_key_prefix.environment.var}"}}, "name": "test_environment", "setting": map[string]interface{}{"name": "MinSize", "namespace": "aws:autoscaling:asg", "value": "1"}}}}}),
		},
		{
			name: "single-provider",
			controlConfigs: []byte(`
provider "aws" {
	version             = "=2.46.0"
	alias                  = "one"
}
`),
			expectedResult: map[string]interface{}(map[string]interface{}{"provider": map[string]interface{}{"aws": map[string]interface{}{"alias": "one", "version": "=2.46.0"}}}),
		},
		{
			name: "two-providers",
			controlConfigs: []byte(`
provider "aws" {
	version             = "=2.46.0"
	alias                  = "one"
}
provider "aws" {
	version             = "=2.47.0"
	alias                  = "two"
}
`),
			expectedResult: map[string]interface{}(map[string]interface{}{"provider": map[string]interface{}{"aws": []interface{}{map[string]interface{}{"alias": "one", "version": "=2.46.0"}, map[string]interface{}{"alias": "two", "version": "=2.47.0"}}}}),
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			var unmarshalledConfigs interface{}

			err := ParseHCL2([]byte(test.controlConfigs), &unmarshalledConfigs)
			assert.Nil(t, err)
			require.NotNil(t, unmarshalledConfigs)
			assert.Equal(t, test.expectedResult, unmarshalledConfigs)
		})
	}
}

func TestYAMLParser(t *testing.T) {
	t.Run("error parsing a YAML document", func(t *testing.T) {

		testTable := []struct {
			name           string
			controlConfigs []byte
			expectedResult interface{}
			shouldError    bool
		}{
			{
				name:           "a single config",
				controlConfigs: []byte(`sample: true`),
				expectedResult: map[string]interface{}{
					"sample": true,
				},
				shouldError: false,
			},
			{
				name: "a single config with multiple yaml subdocs",
				controlConfigs: []byte(`---
sample: true
---
hello: true
---
nice: true`),
				expectedResult: []interface{}{
					map[string]interface{}{
						"sample": true,
					},
					map[string]interface{}{
						"hello": true,
					},
					map[string]interface{}{
						"nice": true,
					},
				},
				shouldError: false,
			},
		}

		for _, test := range testTable {
			t.Run(test.name, func(t *testing.T) {
				var unmarshalledConfigs interface{}

				err := ParseYAML(test.controlConfigs, &unmarshalledConfigs)
				assert.Nil(t, err)
				require.NotNil(t, unmarshalledConfigs)
				assert.Equal(t, test.expectedResult, unmarshalledConfigs)
			})
		}
	})
}
