package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/snyk/snyk-iac-rules/internal"
	"github.com/snyk/snyk-iac-rules/util"
)

var BuildIgnore = append(TestIgnore, "testing", "*_test.rego")

var buildCommand = &cobra.Command{
	Use:   "build [path]",
	Short: "Build an OPA bundle",
	Long: `Build an OPA bundle.

The 'build' command packages OPA policy and data files into bundles. Bundles are
gzipped tarballs. Paths referring to directories are loaded recursively.

To start, run:
$ snyk-iac-rules build
An optional path can be provided if the current directory contains more than just 
the rules for the bundle.

To ignore test files, use the '--ignore' flag:
$ snyk-iac-rules build --ignore testing --ignore "*_test.rego"

If the 'template' command was used to generate the rules, then the default 
entrypoint is "rules/deny". 
Otherwise, you can also override the entrypoint:
$ snyk-iac-rules build --entrypoint "<package name>/<function name>"

The generated bundle has the name 'bundle.tar.gz', but a custom name can be provided:
$ snyk-iac-rules build -o custom-bundle.tar.gz

See our documentation to learn more: 
https://docs.snyk.io/products/snyk-infrastructure-as-code/custom-rules/getting-started-with-the-sdk/bundling-rules
`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("Too many paths provided")
		}
		return nil
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		var path string
		if len(args) == 0 {
			path = "."
		} else {
			path = args[0]
		}
		util.CheckIfRunningInRootDirectory(path)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// add default path if not provided
			args = append(args, ".")
		}
		err := internal.RunBuild(args, buildParams)
		if err != nil {
			return err
		}

		fmt.Printf("Generated bundle: %s \n", buildParams.OutputFile)
		return nil
	},
}

func newBuildCommandParams() *internal.BuildCommandParams {
	return &internal.BuildCommandParams{
		Entrypoint:   util.NewRepeatedStringFlag("rules/deny"),
		Target:       util.NewEnumFlag(internal.TargetWasm, []string{internal.TargetRego, internal.TargetWasm}),
		Capabilities: util.NewCapabilitiesFlag(),
	}
}

var buildParams = newBuildCommandParams()

func init() {
	buildCommand.Flags().VarP(&buildParams.Entrypoint, "entrypoint", "e", "set slash separated entrypoint path")
	buildCommand.Flags().StringVarP(&buildParams.OutputFile, "output", "o", "bundle.tar.gz", "set the output filename")
	buildCommand.Flags().StringSliceVarP(&buildParams.Ignore, "ignore", "", BuildIgnore, "set file and directory names to ignore during loading")
	buildCommand.Flags().VarP(&buildParams.Target, "target", "t", "set the output bundle target type")
	buildCommand.Flags().VarP(&buildParams.Capabilities, "capabilities", "c", "set configurable set of OPA capabilities")

	RootCommand.AddCommand(buildCommand)
}
