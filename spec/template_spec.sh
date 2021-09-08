#!/bin/bash

cleanup() { rm -rf ./fixtures/custom-rules/rules/test; }
AfterAll 'cleanup'

Describe './snyk-iac-custom-rules template ./fixtures/custom-rules --rule test'
   It 'returns passing test status'
      When call ./snyk-iac-custom-rules template ./fixtures/custom-rules --rule test
      The status should be success
      The output should include 'Templating rule...'
      The output should include 'Templated directory'
      The output should include '/fixtures/custom-rules/rules'
      The output should include 'Templated directory'
      The output should include '/fixtures/custom-rules/rules/test'
      The output should include 'Templated file'
      The output should include '/fixtures/custom-rules/rules/test/main.rego'
      The output should include 'Templated file'
      The output should include '/fixtures/custom-rules/rules/test/main_test.rego'
   End
End

Describe './snyk-iac-custom-rules template ./fixtures/custom-rules --rule test'
   It 'returns passing test status'
      When call ./snyk-iac-custom-rules template ./fixtures/custom-rules --rule test
      The status should be failure
      The output should include 'Templating rule...'
      The output should include 'Rule with the provided name already exists'
      The stderr should be present
   End
End