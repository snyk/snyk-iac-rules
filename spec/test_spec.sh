#!/bin/bash
Describe 'go run main.go test ./fixtures/custom-rules'
   It 'returns passing test status'
      When call go run main.go test ./fixtures/custom-rules
      The status should be success
      The output should include 'Executing Rego test cases...
PASS: 3/3'
   End
End
