package rules

deny[msg] {
	input.spec.template.todo
	msg := {
		"publicId": "{{.RuleName}}",
		"title": "<TODO>",
		"severity": "<TODO>",
		"issue": "",
		"impact": "",
		"remediation": "",
		"msg": "spec.template.todo",
		"references": [],
	}
}