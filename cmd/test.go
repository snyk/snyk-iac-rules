package cmd

import (
	"fmt"
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
	Use:   "test [path...]",
	Short: "Execute Rego test cases",
	Long: `Execute Rego test cases.

The 'test' command executes all test cases discovered in matching files. 
Test cases are rules whose names have the prefix "test_".

To start, run:
$ snyk-iac-rules test
An optional list of paths can be provided if the current directory contains more than just 
the rules for the bundle.

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
	PreRunE: func(Cmd *cobra.Command, args []string) error {
		// If an --explain flag was set, turn on verbose output
		if testParams.Explain.IsSet() {
			testParams.Verbose = true
		}

		return nil
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// add default paths if none were provided
			args = append(args, "rules/", "lib/")
		}
		err := util.IsPointingAtTemplatedRules(args)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// add default paths if none were provided
			args = append(args, "rules/", "lib/")
		}
		return internal.RunTest(args, testParams)
	},
}

func newTestCommandParams() *internal.TestCommandParams {
	return &internal.TestCommandParams{
		Explain: util.NewEnumFlag(util.ExplainModeFails, []string{util.ExplainModeFails, util.ExplainModeFull, util.ExplainModeNotes}),
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
