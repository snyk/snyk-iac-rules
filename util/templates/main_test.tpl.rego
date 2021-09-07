package rules

import data.lib
import data.lib.testing

test_{{ call .Replace .RuleName "-" "_" }} {
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

	testing.evaluate_test_cases("{{.RuleName}}", test_cases)
}