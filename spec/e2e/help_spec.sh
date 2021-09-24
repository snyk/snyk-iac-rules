#!/bin/bash
Describe './snyk-iac-rules'
  It 'returns help info'
    When call ./snyk-iac-rules
    The status should be success
    The output should include 'SDK to write, debug, test, and bundle custom rules for Snyk Infrastructure as Code.

Not sure where to start?

1. Run the following command to learn how to generate a scaffolded rule:
$ snyk-iac-rules template --help

2. Run the following command to learn how to parse a file into the JSON structure that Rego understands:
$ snyk-iac-rules parse --help

3. Run the following command to learn how to test a Rego rule:
$ snyk-iac-rules test --help

4. Run the following command to learn how to build the bundle for the Snyk IaC CLI:
$ snyk-iac-rules build --help

Usage:
  snyk-iac-rules [command]

Available Commands:
  build       Build an OPA bundle
  help        Help about any command
  parse       Parse a fixture into JSON format
  template    Template a new rule
  test        Execute Rego test cases

Flags:
  -h, --help   help for snyk-iac-rules

Use "snyk-iac-rules [command] --help" for more information about a command.'
  End
End
