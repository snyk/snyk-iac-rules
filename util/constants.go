package util

const (
	TargetRego = "rego"
	TargetWasm = "wasm"
)

const (
	JSON           = "json"
	HCL2           = "hcl2"
	YAML           = "yaml"
	TERRAFORM_PLAN = "tf-plan"
)

const (
	LOW      = "low"
	MEDIUM   = "medium"
	HIGH     = "high"
	CRITICAL = "critical"
)

var ValidSeverityLevels = []string{LOW, MEDIUM, HIGH, CRITICAL}

const (
	ExplainModeFull  = "full"
	ExplainModeNotes = "notes"
	ExplainModeFails = "fails"
)
