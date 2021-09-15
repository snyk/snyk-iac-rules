package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/snyk/snyk-iac-custom-rules/internal"
)

var templateCommand = &cobra.Command{
	Use:   "template <path>",
	Short: "Template a new rule",
	Long: `Template a new rule.

The 'template' command generates the scaffolding for writing new rules.

To start, run the following command, replacing <rule> with the name of a rule.
$ snyk-iac-rules template --rule <rule>

A rules/ folder is created, which will contain a folder named after the provided
rule name. In this folder the rule definition can be found in the main.rego file,
along with a test in the main_test.rego file.

Each rule must return a JSON structure containing the following required fields:
- 'publicId': the name of the rule (automatically filled in by the command)
- 'title': rule identifier in the snyk CLI
- 'severity': the severity of the rule; can be one of 'low', 'medium', 'high', and 'critical'
- 'msg': the misconfigured path in the fixture

By default, the deny function is where the rule will be implemented. A different 
folder structure can be used, but overrides must be provided when generating the
bundle. To find out how, run the following command:
$ snyk-iac-rules build --help

On top of generating the scaffolding for the rule, the command also generates a
testing framework.
To learn more about this, run the following command:
$ snyk-iac-rules test --help 
`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("Too many paths provided")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// make sure rule name doesn't have any whitespace in it
		if strings.Contains(templateParams.Rule, " ") {
			return fmt.Errorf("Rule name cannot contain whitespace")
		}

		// make sure rule name doesn't belong to Snyk namesapce
		if strings.HasPrefix(templateParams.Rule, "SNYK-") {
			return fmt.Errorf("Rule name cannot start with \"SNYK-\"")
		}

		// prepare directory for templating
		currentDirectory, err := os.Getwd()
		if err != nil {
			return err
		}
		if len(args) == 0 {
			args = append(args, currentDirectory)
		} else {
			args = []string{path.Join(currentDirectory, args[0])}
		}

		err = internal.RunTemplate(args, templateParams)
		if err != nil {
			return err
		}

		fmt.Println("Generated template")
		return nil
	},
}

func newTemplateCommandParams() *internal.TemplateCommandParams {
	return &internal.TemplateCommandParams{}
}

var templateParams = newTemplateCommandParams()

func init() {
	templateCommand.Flags().StringVarP(&templateParams.Rule, "rule", "r", "", "provide rule name")
	templateCommand.MarkFlagRequired("rule") //nolint
	RootCommand.AddCommand(templateCommand)
}
