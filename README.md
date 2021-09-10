# Snyk IaC Custom Rules

[![CircleCI](https://circleci.com/gh/snyk/snyk-iac-custom-rules/tree/develop.svg?style=svg&circle-token=5597b9f0189554f754f38400cbe9d8f8b334c72a)](https://circleci.com/gh/snyk/snyk-iac-custom-rules/tree/develop) [![Shellspec Tests](https://github.com/snyk/snyk-iac-custom-rules/actions/workflows/main.yml/badge.svg)](https://github.com/snyk/snyk-iac-custom-rules/actions/workflows/main.yml)

This is a Golang SDK that will provide flags for writing, debugging, testing, and bundling a custom rules for the Snyk IaC CLI.


## Folder structure
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
│   
└───util - other utility functions used throughout the code
```

## Usage

* See [DEVELOPMENT.md](DEVELOPMENT.md) for how to setup the environment and add a new command.
* See [RELEASE.md](RELEASE.md) for how to release a new version of the SDK.

### Running Locally

1. Clone the repository
2. Build and run the binary: `go build -o snyk-iac-custom-rules .`. Or, alternatively, `go run main.go {command}`.
3. Run the command:
```
$ ./snyk-iac-custom-rules {command}
```

### Testing

Make sure to build the Golang binary first by following the instructions above.
From the project's root folder, run `shellspec` to run [shellspec](https://github.com/shellspec/shellspec) tests.

### Formatting & Linting

To format all files in the current directory and subdirectories, run `go fmt ./...` from the root directory.

To run the linter, run `golangci-lint run -v --timeout 10m` from the root directory.