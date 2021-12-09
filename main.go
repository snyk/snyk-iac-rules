package main

import (
	"os"

	"github.com/snyk/snyk-iac-rules/builtins"
	"github.com/snyk/snyk-iac-rules/cmd"
)

func main() {
	builtins.RegisterHCLBuiltin()
	builtins.RegisterYAMLBuiltin()
	builtins.RegisterTerraformPlanBuiltin()

	if err := cmd.RootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
