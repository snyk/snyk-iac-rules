#!/bin/bash

Describe './snyk-iac-rules test --help'
   It 'returns passing test status'
      When call ./snyk-iac-rules test --help
      The status should be success
      The output should include 'Usage:
  snyk-iac-rules test [path] [flags]

Flags:
      --explain {fails,full,notes}   enable query explanations (default fails)
  -h, --help                         help for test
      --ignore strings               set file and directory names to ignore during loading (default [.*,fixtures])
  -r, --run string                   run only test cases matching the regular expression
      --timeout duration             set test timeout (default 5s)
  -v, --verbose                      set verbose logging mode'
   End
End

Describe './snyk-iac-rules test ./fixtures/custom-rules'
   It 'returns passing test status'
      When call ./snyk-iac-rules test ./fixtures/custom-rules
      The status should be success
      The output should include 'PASS: 3/3'
   End
End

Describe './snyk-iac-rules test ./fixtures/custom-rules --run test_CUSTOM_1'
   It 'returns passing test status'
      When call ./snyk-iac-rules test ./fixtures/custom-rules --run test_CUSTOM_1
      The status should be success
      The output should include 'PASS: 1/1'
   End
End
