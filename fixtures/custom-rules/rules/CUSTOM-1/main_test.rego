package rules

import data.lib
import data.lib.testing

test_CUSTOM_1 {
	test_cases := [{
		"want_msgs": ["input.resource.aws_security_group[denied].tags"],
		"fixture": "test.tf",
	}]

	testing.evaluate_test_cases("CUSTOM-1", "./fixtures/custom-rules/rules/CUSTOM-1/fixtures", test_cases)
}
