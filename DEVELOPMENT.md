# Development

## Tools

### Go

This repo is written in the [Go](https://golang.org) programming language, which can be installed from their [installation page](https://golang.org/doc/install). The version it was developed with was v1.17, which is the minimum required version for this SDK.

### Shellspec
The SDK developed in this repo is tested by [shellspec](https://github.com/shellspec/shellspec), which can be installed as documented in [their README](https://github.com/shellspec/shellspec#installation).

### Golangci-lint
The SDK is linted by [golangci-lint](https://github.com/golangci/golangci-lint), which can be installed as documented in their [local installation page](https://golangci-lint.run/usage/install/#local-installation).

## How to add a new command

The folder structure in this repository is:
```
│   
└───builtins - rego builtins for custom functionality
│
└───cmd - commands and subcommands to register with the cobra CLI  
│   root.go - the root command which needs each subcommand to be registered to
│
└───internal - internal implementation of OPA related functionality
│   
└───fixtures - test fixtures
│   
└───scripts - scripts for CircleCI or GitHub action
│   
└───spec - shellspec tests
│   └─── e2e - shellspec end-to-end tests
|   |
│   └─── contract - shellspec contract tests
│   
└───util - other utility functions used throughout the code
```

To add a new command, you will need to add new files in various of the folders mentioned above so follow the steps below:
1. Add a new file with the name of the command under `cmd`
2. The minimal file should contain:
```go
package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/snyk/snyk-iac-rules/internal"
	"github.com/snyk/snyk-iac-rules/util"
)

var <command>Command = &cobra.Command{
	Use:   "<command>",
	Short: "",
	Long: ``,
	SilenceUsage: true, // disables help from being printed if command fails
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.Run<Command>(args, <command>Params)
	},
}

func new<Command>CommandParams() *internal.<Command>CommandParams {
	return &internal.<Command>CommandParams{
	}
}

var <command>Params = new<Command>CommandParams()

func init() {
	// initialise flags for the command
	RootCommand.AddCommand(<command>Command)
}
```

3. Add a new file with the name of the command under `internal` - this will contain the implementation for `Run<Command>`

```go
package internal

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type <Command>CommandParams struct {
}

func Run<Command>(args []string, params *<Command>CommandParams) error {
	return nil
}
```

4. Add a basic test to with the name of the command under `spec`:

```go
#!/bin/bash
Describe './snyk-iac-rules <command>'
   It 'returns passing test status'
      When call ./snyk-iac-rules <command>
      The status should be success
      The output should include ''
   End
End
```  

### How to run the SDK locally

1. Clone the repository
2. Build and run the binary: `go build -o snyk-iac-rules .`. Or, alternatively, `go run main.go {command}`.
3. Run the command:
```
$ ./snyk-iac-rules {command}
```

### How to run the tests

Make sure to build the Golang binary first by following the instructions in the README.

#### End-to-end tests
**Note** These instructions are for internal use only, as the  feature is behind a feature flag and entitlement.

To run `shellspec` tests, set the following environment variables based on the registry you want to push your bundle to:
```
export OS=<your OS>
export OCI_REGISTRY_NAME=<e.g. docker.io/<username>/<repo>:latest>
export OCI_REGISTRY_USERNAME=<username>
export OCI_REGISTRY_PASSWORD=<password>
```

From the project's root folder, run `shellspec "spec/e2e"` to run [shellspec](https://github.com/shellspec/shellspec) tests for the end-to-end flow. This will verify the behaviour of the commands against the fixtures in the `./fixtures` folder and make sure the output and exit codes are correct.

These tests run as part of the entire CI/CD for our three supported OS's (Linux, MacOs, and Windows), apart from when we release a new version of the SDK.

#### Contract tests with the Snyk CLI
**Note** These instructions are for internal use only, as the  feature is behind a feature flag and entitlement.

To run the contract tests with Snyk, install Snyk by running `npm i -g snyk` and set the `SNYK_TOKEN` to the Auth Token from https://app.snyk.io/account. 

Finally, run `shellspec "spec/contract"` to verify if the generated bundle from the SDK is valid for the Snyk CLI. 

These tests run as part of the entire CI/CD for our three supported OS's (Linux, MacOs, and Windows), apart from when we release a new version of the SDK.

#### End-to-End tests for Supported Registries
**Note** These instructions are for internal use only, as the  feature is behind a feature flag and entitlement.

To run the tests for supported registries, follow the same steps as the contract tests but also set the following environment variables and log into the registries manually:
```
export OCI_DOCKERHUB_REGISTRY_NAME=docker.io/<repo name>:<repo tag>
export OCI_DOCKERHUB_REGISTRY_URL=https://registry-1.${OCI_DOCKERHUB_REGISTRY_NAME}
export OCI_DOCKERHUB_REGISTRY_USERNAME=
export OCI_DOCKERHUB_REGISTRY_PASSWORD=

export OCI_AZURE_REGISTRY_NAME=<registry>/<account name>/<repo name>:<repo tag>
export OCI_AZURE_REGISTRY_URL=https://${OCI_AZURE_REGISTRY_NAME}
export OCI_AZURE_REGISTRY_USERNAME=<client id>
export OCI_AZURE_REGISTRY_PASSWORD=<client secret>

export OCI_HARBOR_REGISTRY_NAME=<registry>/<account name>/<repo name>:<repo tag>
export OCI_HARBOR_REGISTRY_URL=https://${OCI_HARBOR_REGISTRY_NAME}
export OCI_HARBOR_REGISTRY_USERNAME=
export OCI_HARBOR_REGISTRY_PASSWORD=

export OCI_GITHUB_REGISTRY_NAME=<registry>/<repo name>:<repo tag>
export OCI_ELASTIC_REGISTRY_URL=https://${OCI_GITHUB_REGISTRY_NAME}
export OCI_GITHUB_REGISTRY_USERNAME=<GitHub username>
export OCI_GITHUB_REGISTRY_PASSWORD=<personal access token>

export OCI_ELASTIC_REGISTRY_NAME=<registry>/<repo name>:<repo tag>
export OCI_ELASTIC_REGISTRY_URL=https://${OCI_ELASTIC_REGISTRY_NAME}
export OCI_ELASTIC_REGISTRY_USERNAME=AWS
export OCI_ELASTIC_REGISTRY_PASSWORD=$(aws ecr get-login-password)
```

Finally, run `shellspec "spec/registries"` to verify if the generated bundle from the SDK is valid for the Snyk CLI.

These tests run as part of the CI/CD for the `develop` branch on Linux only.

#### Unit tests

Alternatively, you can run `go test ./...` for Golang unit tests. These test the files under `./internal` and `./util` and make sure complex behaviour is maintained between code changes.

For test coverage, run `go test -coverprofile cover.out fmt ./...` and then `go tool cover -html=cover.out` to see what is missing.

These tests run as part of the entire CI/CD, apart from when we release a new version of the SDK.

**Note**
All the fixtures under the `./fixture` folder are used for all Shellspec tests.

### How to lint the code

To format all files in the current directory and subdirectories, run `go fmt ./...` from the root directory.

To run the linter, run `golangci-lint run -v --timeout 10m` from the root directory.
