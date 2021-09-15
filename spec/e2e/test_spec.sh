#!/bin/bash
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
