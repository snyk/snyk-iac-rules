# Snyk IaC Custom Rules

[![CircleCI](https://circleci.com/gh/snyk/snyk-iac-custom-rules/tree/develop.svg?style=svg)](https://circleci.com/gh/snyk/snyk-iac-custom-rules/tree/develop)

This is a Golang CLI that will provide flags for writing, debugging, testing, and bundling a customer's custom rules for the Snyk IaC CLI.

## Usage

### Running Locally

Environment preparation
* Install [Go](https://golang.org/doc/install)
* VSCode - Extentions - [Go](https://marketplace.visualstudio.com/items?itemName=ms-vscode.Go), [Open Policy Agent](https://marketplace.visualstudio.com/items?itemName=tsandall.opa)

1. Clone the repository
2. Build and run the binary: `go build -o synk-iac-custom-rules .`. Or, alternatively, `go run main.go`.
3. Run the command:
```
$ ./synk-iac-custom-rules
```

### Testing

From the project's root folder, run `shellspec` to run [shellspec](https://github.com/shellspec/shellspec) tests.