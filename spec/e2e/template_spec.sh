#!/bin/bash

cleanup() { rm -rf ./fixtures/custom-rules/rules/test; }
AfterAll 'cleanup'

Describe './snyk-iac-rules template --help'
   It 'returns passing test status'
      When call ./snyk-iac-rules template --help
      The status should be success
      The output should include 'Usage:
  snyk-iac-rules template [path] [flags]

Flags:
  -h, --help                                  help for template
  -r, --rule string                           provide rule id
  -s, --severity {low,medium,high,critical}   provide rule severity (default low)
  -t, --title string                          provide rule title (default "Default title")'
   End
End

Describe './snyk-iac-rules template ./fixtures/custom-rules --rule test'
   It 'returns passing test status'
      When call ./snyk-iac-rules template ./fixtures/custom-rules --rule test
      The status should be success
      The output should include 'Template rules directory'
      The output should include 'Template rules/test directory'
      The output should include 'Template rules/test/main.rego file'
      The output should include 'Template rules/test/main_test.rego file'
      The output should include 'Template rules/test/fixtures directory'
      The output should include 'Template rules/test/fixtures/allowed.tf file'
      The output should include 'Template rules/test/fixtures/denied.tf file'
      The output should include 'Generated template'
   End
End

Describe './snyk-iac-rules template ./fixtures/custom-rules --rule test'
   It 'returns passing test status'
      When call ./snyk-iac-rules template ./fixtures/custom-rules --rule test
      Dump
      The status should be failure
      The output should include 'Rule with the provided name already exists'
      The stderr should be present
   End
End
