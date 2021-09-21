package rules

import data.lib
import data.lib.testing

test_{{ call .Replace .RuleID "-" "_" }} {
	# array containing test cases where the rule is allowed
	allowed_test_cases := [{
		"want_msgs": [],
		"fixture": "allowed.json",
	}]

	# array containing cases where the rule is denied
	denied_test_cases := [{
		"want_msgs": ["input.resource.test[denied].todo"], # verifies that the correct msg is returned by the denied rule
		"fixture": "denied.json",
	}]

	test_cases := array.concat(allowed_test_cases, denied_test_cases)
	testing.evaluate_test_cases("{{.RuleID}}", test_cases)
}
