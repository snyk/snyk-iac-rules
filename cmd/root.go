package cmd

import (
	"github.com/spf13/cobra"
)

var version string = "dev"

// RootCommand is the base CLI command that all subcommands are added to.
var RootCommand = NewRootCmd()

func NewRootCmd() cobra.Command {
	rootCommand := cobra.Command{
		Use:   "snyk-iac-rules",
		Short: "Snyk IaC Custom Rules",
		Long: `SDK to write, debug, test, and bundle custom rules for Snyk Infrastructure as Code.

Not sure where to start?

1. Run the following command to learn how to generate a scaffolded rule:
$ snyk-iac-rules template --help

2. Run the following command to learn how to parse a file into the JSON structure that Rego understands:
$ snyk-iac-rules parse --help

3. Run the following command to learn how to test a Rego rule:
$ snyk-iac-rules test --help

4. Run the following command to learn how to build the bundle for the Snyk IaC CLI:
$ snyk-iac-rules build --help

5. Verify your bundle locally using the Snyk IaC CLI:
$ snyk iac test --rules --help

6. Run the following command to learn how to push the bundle to an OCI Registry:
$ snyk-iac-rules push --help

7. Use the bundle stored in the OCI registry in the Snyk IaC CLI:
$ snyk --help

See our documentation to learn more:
https://docs.snyk.io/products/snyk-infrastructure-as-code/custom-rules
`,
		Version: version,
	}
	rootCommand.CompletionOptions.DisableDefaultCmd = true
	return rootCommand
}
