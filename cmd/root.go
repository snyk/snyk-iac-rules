package cmd

import (
	"github.com/spf13/cobra"
)

// RootCommand is the base CLI command that all subcommands are added to.
var RootCommand = NewRootCmd()

func NewRootCmd() cobra.Command {
	rootCommand := cobra.Command{
		Use:   "snyk-iac-custom-rules",
		Short: "Snyk IaC Custom Rules",
		Long:  "An SDK to write, debug, test, and bundle custom rules for Snyk IaC.",
	}
	rootCommand.CompletionOptions.DisableDefaultCmd = true
	return rootCommand
}
