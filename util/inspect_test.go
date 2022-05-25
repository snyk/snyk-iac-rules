package util

import (
	"fmt"
	"testing"

	"github.com/open-policy-agent/opa/util/test"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveRulesWitInvalidPath(t *testing.T) {
	_, err := RetrieveRules([]string{"./invalid"})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "./invalid: no such file or directory")
}

func TestRetrieveRulesWithNoFiles(t *testing.T) {
	files := map[string]string{}
	test.WithTempFS(files, func(root string) {
		rules, err := RetrieveRules([]string{root})
		assert.Nil(t, err)
		assert.Equal(t, 0, len(rules))
	})
}

func TestRetrieveRulesWithNoRules(t *testing.T) {
	files := map[string]string{
		"test.json": `
			test
		`,
	}

	test.WithTempFS(files, func(root string) {
		rules, err := RetrieveRules([]string{root})
		assert.Nil(t, err)
		assert.Equal(t, 0, len(rules))
	})
}

func TestRetrieveRulesWithRulesWithPublicIdAndSeverityLevel(t *testing.T) {
	files := map[string]string{
		"test.rego": `
			package test
			msg = {
				"publicId": "1",
				"severity": "low"
			}
		`,
	}

	test.WithTempFS(files, func(root string) {
		rules, err := RetrieveRules([]string{root})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(rules))
		assert.Equal(t, "1", rules[0].PublicId)
		assert.Equal(t, "low", rules[0].SeverityLevel)
	})
}

func TestRetrieveRulesWithRulesWithoutPublicId(t *testing.T) {
	files := map[string]string{
		"test.rego": `
			package test
			msg = {
				"publicI":
					"1"
			}
		`,
	}

	test.WithTempFS(files, func(root string) {
		rules, err := RetrieveRules([]string{root})
		assert.Nil(t, err)
		assert.Equal(t, 0, len(rules))
	})
}

func TestRetrieveRulesWithRulesWithDistinctPublicIds(t *testing.T) {
	files := map[string]string{
		"test1.rego": `
			package test
			msg = {
				"publicId":
					"1"
			}
		`,
		"test2.rego": `
			package test
			msg = {
				"publicId":
					"2"
			}
		`,
		"test2_test.rego": `
			package test
			msg = {
				"publicId":
					"3"
			}
		`,
	}

	test.WithTempFS(files, func(root string) {
		rules, err := RetrieveRules([]string{root})
		assert.Nil(t, err)
		assert.Equal(t, 2, len(rules))
		assert.Subset(t, rules, []Rule{
			{
				PublicId: "1",
				Path:     fmt.Sprintf("%s/test1.rego", root),
			},
			{
				PublicId: "2",
				Path:     fmt.Sprintf("%s/test2.rego", root),
			},
		})
	})
}

func TestRetrieveRulesWithRulesWithSamePublicIds(t *testing.T) {
	files := map[string]string{
		"test1.rego": `
			package test
			msg = {
				"publicId":
					"1"
			}
		`,
		"test2.rego": `
			package test
			msg = {
				"publicId":
					"1"
			}
		`,
		"test2_test.rego": `
			package test
			msg = {
				"publicId":
					"3"
			}
		`,
	}

	test.WithTempFS(files, func(root string) {
		rules, err := RetrieveRules([]string{root})
		assert.Nil(t, err)
		assert.Equal(t, 2, len(rules))
		assert.Subset(t, rules, []Rule{
			{
				PublicId: "1",
				Path:     fmt.Sprintf("%s/test1.rego", root),
			},
			{
				PublicId: "1",
				Path:     fmt.Sprintf("%s/test2.rego", root),
			},
		})
	})
}
