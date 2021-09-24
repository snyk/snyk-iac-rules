#!/bin/bash

cleanup() { rm -rf ./fixtures/custom-rules/rules/test; }
AfterAll 'cleanup'

Describe './snyk-iac-rules template ./fixtures/custom-rules --rule test'
   It 'returns passing test status'
      When call ./snyk-iac-rules template ./fixtures/custom-rules --rule test
      The status should be success
      The stderr should include 'Template rules directory'
      The stderr should include 'Template rules/test directory'
      The stderr should include 'Template rules/test/main.rego file'
      The stderr should include 'Template rules/test/main_test.rego file'
      The stderr should include 'Template rules/test/fixtures directory'
      The stderr should include 'Template rules/test/fixtures/allowed.tf file'
      The stderr should include 'Template rules/test/fixtures/denied.tf file'
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
