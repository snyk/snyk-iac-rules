package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/snyk/snyk-iac-rules/internal"
	"github.com/snyk/snyk-iac-rules/util"
)

var templateCommand = &cobra.Command{
	Use:   "template [path]",
	Short: "Template a new rule",
	Long: `Generate the template for a new rule.

The 'template' command generates the scaffolding for writing new rules.

To start, run the following command, replacing <rule> with the name of a rule and <format> with one of the formats.
$ snyk-iac-rules template --rule <rule> --format <format>

A rules/ folder is created, which will contain a folder named after the provided
rule name. In this folder the rule definition can be found in the 'main.rego' file,
along with a test in the 'main_test.rego' file.

Each rule must return a JSON structure containing the following required fields:
- 'publicId': the name of the rule (automatically filled in by the command)
- 'title': rule identifier in the snyk CLI
- 'severity': the severity of the rule; can be one of 'low', 'medium', 'high', and 'critical'
- 'msg': the misconfigured path in the fixture

See our documentation to learn more: 
https://docs.snyk.io/products/snyk-infrastructure-as-code/custom-rules/getting-started-with-the-sdk/writing-a-rule

By default, the deny function is where the rule will be implemented. A different 
folder structure can be used, but overrides must be provided when generating the
bundle. To find out how, run the following command:
$ snyk-iac-rules build --help

On top of generating the scaffolding for the rule, the command also generates a
testing framework. To learn more about this, run the following command:
$ snyk-iac-rules test --help 
`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("Too many paths provided")
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		ruleId := strings.ToUpper(templateParams.RuleID)

		// make sure rule name doesn't have any whitespace in it
		if strings.Contains(ruleId, " ") {
			return fmt.Errorf("Rule name cannot contain whitespace")
		}

		// make sure rule name doesn't belong to Snyk namespace
		if strings.HasPrefix(ruleId, "SNYK-") {
			return fmt.Errorf("Rule name cannot start with \"SNYK-\"")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// add default path if not provided
			args = append(args, ".")
		}
		err := internal.RunTemplate(args, templateParams)
		if err != nil {
			return err
		}

		fmt.Println("Generated template")
		return nil
	},
}

func newTemplateCommandParams() *internal.TemplateCommandParams {
	return &internal.TemplateCommandParams{
		RuleSeverity: util.NewEnumFlag(internal.LOW, []string{internal.LOW, internal.MEDIUM, internal.HIGH, internal.CRITICAL}),
		RuleFormat:   util.NewEnumFlag("", []string{internal.HCL2, internal.JSON, internal.YAML, internal.TERRAFORM_PLAN}),
	}
}

var templateParams = newTemplateCommandParams()

func init() {
	templateCommand.Flags().StringVarP(&templateParams.RuleID, "rule", "r", "", "provide rule id")
	templateCommand.Flags().VarP(&templateParams.RuleFormat, "format", "f", "provide rule format")
	templateCommand.Flags().StringVarP(&templateParams.RuleTitle, "title", "t", "Default title", "provide rule title")
	templateCommand.Flags().VarP(&templateParams.RuleSeverity, "severity", "s", "provide rule severity")
	templateCommand.MarkFlagRequired("rule")   //nolint
	templateCommand.MarkFlagRequired("format") //nolint
	RootCommand.AddCommand(templateCommand)
}
