package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/snyk/snyk-iac-custom-rules/internal"
)

var templateCommand = &cobra.Command{
	Use:   "template <path>",
	Short: "Template a new rule",
	Long: `Template a new rule.

The 'template' command generates some templating code for writing new rules.
`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Templating rule...")
		if len(args) == 0 {
			args = append(args, "./")
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
	return &internal.TemplateCommandParams{}
}

var templateParams = newTemplateCommandParams()

func init() {
	templateCommand.Flags().StringVarP(&templateParams.Rule, "rule", "r", "", "provide rule name")
	templateCommand.MarkFlagRequired("rule") //nolint
	RootCommand.AddCommand(templateCommand)
}
