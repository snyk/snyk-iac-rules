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

To start, log into your OCI registry locally and then run:
$ snyk-iac-rules push <path to generated bundle> --registry <desired location of OCI artifact in your OCI registry>
If the tag is not provided at the end of the provided registry, the tool defaults to 'latest'.

For example, if your generated bundle is called 'bundle.tar.gz' and the OCI registry you're using is DockerHub, the following command must be run: 
$ snyk-iac-rules push bundle.tar.gz --registry docker.io/example/bundle

Once the bundle has been pushed to an OCI registry, the Snyk IaC CLI can be configured to pull it down for scanning:
$ snyk config set oci-registry-url=https://registry-1.docker.io/example/bundle:latest
$ snyk config set oci-registry-username=username
$ snyk config set oci-registry-password=password

See our documentation to learn more: 
https://docs.snyk.io/products/snyk-infrastructure-as-code/custom-rules/getting-started-with-the-sdk/pushing-a-bundle
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

		if strings.Contains(pushParams.BundleRegistry, "://") {
			return fmt.Errorf("The provided container registry includes a protocol. Please remove it and try again")
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
