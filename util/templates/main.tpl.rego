package rules

deny[msg] {
	resource := input.resource.test[name]
	resource.todo
	msg := {
		"publicId": "{{.RuleID}}",
		"title": "{{.RuleTitle}}",
		"severity": "{{.RuleSeverity}}",
		"issue": "",
		"impact": "",
		"remediation": "",
		"msg": sprintf("input.resource.test[%s].todo", [name]),
		"references": [],
	}
}
