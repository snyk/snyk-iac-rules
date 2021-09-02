package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/snyk/snyk-iac-custom-rules/internal"
	"github.com/snyk/snyk-iac-custom-rules/util"
)

var testCommand = &cobra.Command{
	Use:   "test <path>",
	Short: "Execute Rego test cases",
	Long: `Execute Rego test cases.

The 'test' command takes a file or directory path as input and executes all
test cases discovered in matching files. Test cases are rules whose names have the prefix "test_".
`,
	SilenceUsage: true,
	PreRunE: func(Cmd *cobra.Command, args []string) error {
		// If an --explain flag was set, turn on verbose output
		if testParams.Explain.IsSet() {
			testParams.Verbose = true
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Executing Rego test cases...")
		if len(args) == 0 {
			args = append(args, "./")
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
	testCommand.Flags().StringSliceVarP(&testParams.Ignore, "ignore", "", []string{".*", "scripts", "build"}, "set file and directory names to ignore during loading (e.g., '.*' excludes hidden files)")
	RootCommand.AddCommand(testCommand)
}
