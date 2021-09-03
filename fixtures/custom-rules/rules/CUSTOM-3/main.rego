package rules

import data.lib

cidr_block_contains_all_networking(cidr_blocks) {
	contains(cidr_blocks[_], "0.0.0.0/0")
}

cidr_block_contains_all_networking(cidr_blocks) {
	contains(cidr_blocks[_], "::/0")
}

deny[msg] {
	resource := input.resource.aws_security_group[name]
	ingres_rules := lib.normalize_to_array(resource.ingress)
	rule := ingres_rules[i]
	cidr_block_contains_all_networking(rule.cidr_blocks)

	msg := {
		"publicId": "CUSTOM-3",
		"title": "Test",
		"subType": "",
		"severity": "critical",
		"issue": "",
		"impact": "",
		"remediation": "",
		"msg": sprintf("input.resource.aws_security_group[%v].ingress[%v]", [name, i]),
		"references": [],
	}
}
