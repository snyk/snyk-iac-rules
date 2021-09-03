package rules

deny[msg] {
	resource := input.resource.aws_security_group[name]
	not resource.tags
	msg := {
		"publicId": "CUSTOM-1",
		"title": "Missing tags",
		"subType": "",
		"severity": "low",
		"issue": "Missing tags",
		"impact": "Depends",
		"remediation": "Set `aws_security_group.tags`",
		"msg": sprintf("input.resource.aws_security_group[%s].tags", [name]),
		"references": [],
	}
}
