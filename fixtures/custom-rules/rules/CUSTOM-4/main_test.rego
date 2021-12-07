package rules

import data.lib
import data.lib.testing

test_CUSTOM_4 {
	# array containing test cases where the rule is allowed
	allowed_test_cases := []

	# array containing cases where the rule is denied
	denied_test_cases := [{
		"want_msgs": ["input.resource.aws_vpc[example].tags"],
		"fixture": "test.json.tfplan",
	}]

	test_cases := array.concat(allowed_test_cases, denied_test_cases)
	testing.evaluate_test_cases("CUSTOM-4", "./fixtures/custom-rules/rules/CUSTOM-4/fixtures", test_cases)
}
