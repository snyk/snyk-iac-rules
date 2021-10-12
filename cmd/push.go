package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/snyk/snyk-iac-rules/internal"
	"github.com/snyk/snyk-iac-rules/util"
)

var pushCommand = &cobra.Command{
	Use:   "push <path>",
	Short: "Push generated bundle",
	Long: `Push a custom rules bundle to an OCI registry.

To start, run:
$ snyk-iac-rules push <path to generated bundle> --registry <container registry>
If the tag is not provided at the end of the provided registry, the tool defaults to 'latest'.

The 'push' command takes in the container registry and pushes the bundle at the provided 
path. The container registry must support OCI artifacts.
`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Expected a path to be provided via the command")
		}
		if len(args) > 1 {
			return errors.New("Too many paths provided")
		}
		if !strings.HasSuffix(args[0], ".tar.gz") {
			return errors.New("The path must be to a generated .tar.gz bundle")
		}

		_, err := util.ValidateFilePath(args[0])
		if err != nil {
			return err
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !strings.Contains(pushParams.BundleRegistry, "/") {
			return fmt.Errorf("The provided container registry is invalid")
		}

		repository := strings.Split(pushParams.BundleRegistry, "/")
		tag := repository[len(repository)-1]
		if !strings.Contains(tag, ":") {
			pushParams.BundleRegistry = pushParams.BundleRegistry + ":latest"
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := internal.RunPush(args, pushParams)
		if err != nil {
			return err
		}

		fmt.Println("Successfully pushed bundle to " + pushParams.BundleRegistry)
		return nil
	},
}

func newPushCommandParams() *internal.PushCommandParams {
	return &internal.PushCommandParams{}
}

var pushParams = newPushCommandParams()

func init() {
	pushCommand.Flags().StringVarP(&pushParams.BundleRegistry, "registry", "r", "", "provide container registry")
	pushCommand.MarkFlagRequired("registry") //nolint
	RootCommand.AddCommand(pushCommand)
}
