package rules

import data.lib
import data.lib.testing

test_CUSTOM_2 {
	test_cases := [{
		"want_msgs": ["resource.aws_ami[denied].ebs_block_device[0]"],
		"fixture": {"resource": {"aws_ami": {"denied": {
			"ebs_block_device": {
				"device_name": "/dev/xvda",
				"encrypted": false,
				"volume_size": 8,
			},
			"name": "aws_ami_not_encrypted",
			"root_device_name": "/dev/xvda",
			"virtualization_type": "hvm",
		}}}},
	}]

	testing.evaluate_test_cases("CUSTOM-2", test_cases)
}
