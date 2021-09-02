#!/bin/bash
Describe 'go run main.go'
   It 'returns help info'
      When call go run main.go
      The status should be success
      The output should include 'NAME:
   snyk-iac-custom-rules - Use this SDK to write, debug, test, and bundle custom rules for the Snyk IaC CLI

USAGE:'
   The output should include 'COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)'
   End
End
