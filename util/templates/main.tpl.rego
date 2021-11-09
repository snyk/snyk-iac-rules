package rules

deny[msg] {
	resource := input.resource.test[name]
	resource.todo
	msg := {
		# Mandatory fields
		"publicId": "{{.RuleID}}",
		"title": "{{.RuleTitle}}",
		"severity": "{{.RuleSeverity}}",
		"msg": sprintf("input.resource.test[%s].todo", [name]), # must be the JSON path to the resource field that triggered the deny rule
		# Optional fields
		"issue": "",
		"impact": "",
		"remediation": "",
		"references": [],
	}
}
