package main

import (
	"fmt"
	"os"

	"github.com/snyk/snyk-iac-custom-rules/cmd"
)

func main() {
	if err := cmd.RootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
