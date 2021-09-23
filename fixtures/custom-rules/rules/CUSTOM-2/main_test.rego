package rules

import data.lib
import data.lib.testing

test_CUSTOM_2 {
	test_cases := [{
		"want_msgs": ["resource.aws_ami[denied].ebs_block_device[0]"],
		"fixture": "test.tf",
	}]

	testing.evaluate_test_cases("CUSTOM-2", "./fixtures/custom-rules/rules/CUSTOM-2/fixtures", test_cases)
}
