#!/bin/bash
Describe 'go run main.go build ./fixtures/custom-rules' --ignore testing --ignore "*_test.rego"
   It 'returns passing test status'
      When call go run main.go build ./fixtures/custom-rules --ignore testing --ignore "*_test.rego"
      The status should be success
      The output should include 'Building OPA WASM bundle...'
   End
End
