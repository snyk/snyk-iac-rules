package rules

deny[msg] {
	resource := input.resource.test[name]
	resource.todo
	msg := {
		# Mandatory fields
		"publicId": "{{.RuleID}}",
		"title": "{{.RuleTitle}}",
		"severity": "{{.RuleSeverity}}",
		"msg": sprintf("input.resource.test[%s].todo", [name]),
		# Optional fields
		"issue": "",
		"impact": "",
		"remediation": "",
		"references": [],
	}
}
