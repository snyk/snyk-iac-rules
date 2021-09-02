package rules

import data.lib

deny[msg] {
	resource := input.resource.aws_ami[name]

	ebs_block_device_rules := lib.normalize_to_array(resource.ebs_block_device)
	rule := ebs_block_device_rules[i]

	not rule.encrypted == true
	not rule.snapshot_id

	msg := {
		"publicId": "CUSTOM-2",
		"title": "AMI snapshot is not encrypted",
		"subType": "",
		"severity": "medium",
		"issue": "The AMI snapshot is not encrypted",
		"impact": "Data stored in the snapshot may be sensitive. Without encryption the data may be accessed without appropriate authorization",
		"remediation": "Set `ebs_block_device_rules.encrypted` attribute to `true`",
		"msg": sprintf("resource.aws_ami[%v].ebs_block_device[%v]", [name, i]),
		"references": [],
	}
}
