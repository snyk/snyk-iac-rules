#!/bin/bash
Describe 'go run main.go'
   It 'returns help info'
      When call go run main.go
      The status should be success
      The output should include 'An SDK to write, debug, test, and bundle custom rules for Snyk IaC.'
   End
End
