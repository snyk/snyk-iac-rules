#!/bin/bash

cleanup() { rm bundle.tar.gz; }
AfterAll 'cleanup'

Describe './snyk-iac-rules build ./fixtures/custom-rules --ignore testing --ignore "*_test.rego"'
   It 'returns passing test status'
      When call ./snyk-iac-rules build ./fixtures/custom-rules --ignore testing --ignore "*_test.rego"
      The status should be success
      The output should include 'Generated bundle: bundle.tar.gz'
   End
End
