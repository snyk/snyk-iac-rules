package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/snyk/snyk-iac-custom-rules/internal"
	"github.com/snyk/snyk-iac-custom-rules/util"
)

var BuildIgnore = append(TestIgnore, "*_test.rego")

var buildCommand = &cobra.Command{
	Use:   "build <path>",
	Short: "Build an OPA WASM bundle",
	Long: `Build an OPA WASM bundle.

The 'build' command packages OPA policy and data files into gzipped tarballs containing policies and data written as WASM.
`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Building OPA WASM bundle...")
		if len(args) == 0 {
			currentDirectory, err := os.Getwd()
			if err != nil {
				return err
			}
			args = append(args, currentDirectory)
		}
		err := internal.RunBuild(args, buildParams)
		if err != nil {
			return err
		}

		fmt.Printf("Bundle %s has been generated\n", buildParams.OutputFile)
		return nil
	},
}

func newBuildCommandParams() *internal.BuildCommandParams {
	return &internal.BuildCommandParams{
		Entrypoint: util.NewRepeatedStringFlag("rules/deny"),
		Target:     util.NewEnumFlag(internal.TargetWasm, []string{internal.TargetRego, internal.TargetWasm}),
	}
}

var buildParams = newBuildCommandParams()

func init() {
	buildCommand.Flags().VarP(&buildParams.Entrypoint, "entrypoint", "e", "set slash separated entrypoint path")
	buildCommand.Flags().StringVarP(&buildParams.OutputFile, "output", "o", "bundle.tar.gz", "set the output filename")
	buildCommand.Flags().StringSliceVarP(&buildParams.Ignore, "ignore", "", BuildIgnore, "set file and directory names to ignore during loading (e.g., '.*' excludes hidden files)")
	buildCommand.Flags().VarP(&buildParams.Target, "target", "t", "set the output bundle target type")
	RootCommand.AddCommand(buildCommand)
}
