package rules

import data.lib

deny[msg] {
	resource := input.resource.aws_vpc[name]
    not resource.tags.owner

	msg := {
		"publicId": "CUSTOM-4",
		"title": "Test TFPLAN",
		"subType": "","severity": "medium",
		"msg": sprintf("input.resource.aws_vpc[%s].tags", [name]),
		"issue": "",
		"impact": "",
		"remediation": "",
		"references": [],
	}
}
