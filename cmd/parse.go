package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/snyk/snyk-iac-rules/internal"
	"github.com/snyk/snyk-iac-rules/util"
)

var parseCommand = &cobra.Command{
	Use:   "parse <path>",
	Short: "Parse a fixture into JSON format",
	Long: `Parse a fixture into JSON format.

The 'parse' command takes the path to a fixture and returns the JSON format that 
would need to be used when writing the Rego rules.

For example, to parse a Terraform file run the following command:
$ snyk-iac-rules parse test.tf --format hcl2
The '--format' flag can be left out when parsing Terraform files, as we default to hcl2.

To parse a YAML file run the following command:
$ snyk-iac-rules parse test.yaml --format yaml

And to parse a Terraform plan output in JSON format (Terraform plan output is binary by default, to convert it to JSON you can run 'terraform show -json tfplan.binary > tf-plan.json') run the following command:
$ snyk-iac-rules parse tf-plan.json --format tf-plan

See our documentation to learn more: 
https://docs.snyk.io/products/snyk-infrastructure-as-code/custom-rules/getting-started-with-the-sdk/parsing-an-input-file
`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Expected a path to be provided via the command")
		}
		if len(args) > 1 {
			return errors.New("Too many paths provided")
		}

		fileInfo, err := util.ValidateFilePath(args[0])
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			return errors.New("A path to a directory cannot be provided")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.RunParse(args, parseParams)
	},
}

func newParseCommandParams() *internal.ParseCommandParams {
	return &internal.ParseCommandParams{
		Format: util.NewEnumFlag(util.HCL2, []string{util.HCL2, util.YAML, util.TERRAFORM_PLAN}),
	}
}

var parseParams = newParseCommandParams()

func init() {
	parseCommand.Flags().VarP(&parseParams.Format, "format", "f", "choose the format for the parser")
	RootCommand.AddCommand(parseCommand)
}
