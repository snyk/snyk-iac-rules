package rules

import data.lib
import data.lib.testing

test_CUSTOM_1 {
	test_cases := [{
		"want_msgs": ["input.resource.aws_security_group[denied].tags"],
		"fixture": {"resource": {"aws_security_group": {"denied": {
			"description": "Allow SSH inbound from anywhere",
			"ingress": {
				"cidr_blocks": ["0.0.0.0/0"],
				"from_port": 22,
				"protocol": "tcp",
				"to_port": 22,
			},
			"name": "allow_ssh",
			"vpc_id": "${aws_vpc.main.id}",
		}}}},
	}]

	testing.evaluate_test_cases("CUSTOM-1", test_cases)
}
