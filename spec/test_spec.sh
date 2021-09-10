#!/bin/bash
Describe './snyk-iac-custom-rules test ./fixtures/custom-rules'
   It 'returns passing test status'
      When call ./snyk-iac-custom-rules test ./fixtures/custom-rules
      The status should be success
      The output should include 'Executing Rego test cases...
PASS: 3/3'
   End
End

Describe './snyk-iac-custom-rules test ./fixtures/custom-rules --run test_CUSTOM_1'
   It 'returns passing test status'
      When call ./snyk-iac-custom-rules test ./fixtures/custom-rules --run test_CUSTOM_1
      The status should be success
      The output should include 'Executing Rego test cases...
PASS: 1/1'
   End
End
