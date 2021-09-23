package rules

import data.lib
import data.lib.testing

test_{{ call .Replace .RuleName "-" "_" }} {
	test_cases := [{
		"want_msgs": [],
		"fixture": "allowed.json",
	}, {
		"want_msgs": ["spec.template.todo"],
		"fixture": "denied.json",
	}]

	testing.evaluate_test_cases("{{.RuleName}}", "rules/{{.RuleName}}/fixtures", test_cases)
}
