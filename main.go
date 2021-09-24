package main

import (
	"fmt"
	"os"

	"github.com/snyk/snyk-iac-rules/builtins"
	"github.com/snyk/snyk-iac-rules/cmd"
)

func main() {
	builtins.RegisterHCLBuiltin()
	builtins.RegisterYAMLBuiltin()
	if err := cmd.RootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
