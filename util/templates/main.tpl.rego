package rules

deny[msg] {
	input.spec.template.todo
	msg := {
		"publicId": "{{.RuleID}}",
		"title": "{{.RuleTitle}}",
		"severity": "{{.RuleSeverity}}",
		"issue": "",
		"impact": "",
		"remediation": "",
		"msg": "spec.template.todo",
		"references": [],
	}
}
