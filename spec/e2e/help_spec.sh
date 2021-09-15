#!/bin/bash
Describe './snyk-iac-custom-rules'
   It 'returns help info'
      When call ./snyk-iac-custom-rules
      The status should be success
      The output should include 'An SDK to write, debug, test, and bundle custom rules for Snyk IaC.

Usage:
  snyk-iac-custom-rules [command]

Available Commands:
  build       Build an OPA WASM bundle
  help        Help about any command
  parse       Parse a fixture into JSON format
  template    Template a new rule
  test        Execute Rego test cases

Flags:
  -h, --help   help for snyk-iac-custom-rules

Use "snyk-iac-custom-rules [command] --help" for more information about a command.'
   End
End
