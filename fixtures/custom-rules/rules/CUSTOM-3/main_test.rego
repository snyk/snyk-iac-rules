package rules

import data.lib
import data.lib.testing

test_CUSTOM_3 {
	test_cases := [{
		"want_msgs": ["spec.externalIPs"],
		"fixture": "test.yaml",
	}]

	testing.evaluate_test_cases("CUSTOM-3", "./fixtures/custom-rules/rules/CUSTOM-3/fixtures", test_cases)
}
