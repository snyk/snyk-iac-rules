package cmd

import (
	"errors"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/snyk/snyk-iac-rules/internal"
	"github.com/snyk/snyk-iac-rules/util"
)

var TestIgnore = []string{
	".*",
	"fixtures",
}

var testCommand = &cobra.Command{
	Use:   "test [path]",
	Short: "Execute Rego test cases",
	Long: `Execute Rego test cases.

The 'test' command executes all test cases discovered in matching files. 
Test cases are rules whose names have the prefix "test_".

To start, run:
$ snyk-iac-rules test
An optional path can be provided if the current directory contains more than just the rules for the bundle.

The command can run one test at a time with the help of the --run flag as such:
$ snyk-iac-rules test --run test_<rule>

The testing framework provided as part of this SDK exposes a 'evaluate_test_cases' function from the 'data.lib.testing'
package, which takes the name of a rule as a first argument and an array of test cases structured as such:
- 'want_msgs': array containing the 'msg' field from failing fixtures
- 'fixture': the path to the fixture file in the rule's 'fixture/' folder

See our documentation to learn more: 
https://docs.snyk.io/products/snyk-infrastructure-as-code/custom-rules/getting-started-with-the-sdk/testing-a-rule
`,
	SilenceUsage:  true,
	SilenceErrors: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("Too many paths provided")
		}
		return nil
	},
	PreRunE: func(Cmd *cobra.Command, args []string) error {
		// If an --explain flag was set, turn on verbose output
		if testParams.Explain.IsSet() {
			testParams.Verbose = true
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			currentDirectory, err := os.Getwd()
			if err != nil {
				return err
			}
			args = append(args, currentDirectory)
		}
		return internal.RunTest(args, testParams)
	},
}

func newTestCommandParams() *internal.TestCommandParams {
	return &internal.TestCommandParams{
		Explain: util.NewEnumFlag(internal.ExplainModeFails, []string{internal.ExplainModeFails, internal.ExplainModeFull, internal.ExplainModeNotes}),
	}
}

var testParams = newTestCommandParams()

func init() {
	testCommand.Flags().BoolVarP(&testParams.Verbose, "verbose", "v", false, "set verbose logging mode")
	testCommand.Flags().VarP(&testParams.Explain, "explain", "", "enable query explanations")
	testCommand.Flags().DurationVar(&testParams.Timeout, "timeout", 5*time.Second, "set test timeout")
	testCommand.Flags().StringSliceVarP(&testParams.Ignore, "ignore", "", TestIgnore, "set file and directory names to ignore during loading")
	testCommand.Flags().StringVarP(&testParams.RunRegex, "run", "r", "", "run only test cases matching the regular expression")
	RootCommand.AddCommand(testCommand)
}
