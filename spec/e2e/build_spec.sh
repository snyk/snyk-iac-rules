#!/bin/bash

cleanup() { rm bundle.tar.gz; }
AfterAll 'cleanup'

Describe './snyk-iac-rules build --help'
   It 'returns passing test status'
      When call ./snyk-iac-rules build --help
      The status should be success
      The output should include 'Usage:
  snyk-iac-rules build [path] [flags]

Flags:
  -c, --capabilities string   set configurable set of OPA capabilities
  -e, --entrypoint string     set slash separated entrypoint path (default "rules/deny")
  -h, --help                  help for build
      --ignore strings        set file and directory names to ignore during loading (default [.*,fixtures,testing,*_test.rego])
  -o, --output string         set the output filename (default "bundle.tar.gz")
  -t, --target {rego,wasm}    set the output bundle target type (default wasm)'
   End
End

Describe './snyk-iac-rules build ./fixtures/custom-rules --ignore testing --ignore "*_test.rego"'
   It 'returns passing test status'
      When call ./snyk-iac-rules build ./fixtures/custom-rules --ignore testing --ignore "*_test.rego"
      The status should be success
      The output should include 'Generated bundle: bundle.tar.gz'
   End
End
