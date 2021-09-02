package rules

import data.lib
import data.lib.testing

test_CUSTOM_3 {
	test_cases := [{
		"want_msgs": ["input.resource.aws_security_group[denied].ingress[0]"],
		"fixture": {"resource": {"aws_security_group": {"denied": {
			"description": "Allow TCP inbound from anywhere",
			"ingress": {
				"cidr_blocks": ["::/0"],
				"from_port": 3389,
				"protocol": "tcp",
				"to_port": 3389,
			},
			"name": "allow_tcp",
			"vpc_id": "arn",
		}}}},
	}]

	testing.evaluate_test_cases("CUSTOM-3", test_cases)
}
