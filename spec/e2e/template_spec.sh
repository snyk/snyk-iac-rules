#!/bin/bash

Describe './snyk-iac-rules template --help'
   It 'returns passing test status'
      When call ./snyk-iac-rules template --help
      The status should be success
      The output should include 'Usage:
  snyk-iac-rules template [path] [flags]

Flags:
  -f, --format {hcl2,json,yaml,tf-plan}       provide rule format
  -h, --help                                  help for template
  -r, --rule string                           provide rule id
  -s, --severity {low,medium,high,critical}   provide rule severity (default low)
  -t, --title string                          provide rule title (default "Default title")'
   End
End

Describe './snyk-iac-rules template ./fixtures/custom-rules --format hcl2'
   It 'returns an error'
      When call ./snyk-iac-rules template ./fixtures/custom-rules --format hcl2
      The status should be failure
      The stderr should include 'required flag(s) "rule" not set'
   End
End

Describe './snyk-iac-rules template ./fixtures/custom-rules --rule test'
   It 'returns  an error'
      When call ./snyk-iac-rules template ./fixtures/custom-rules --rule test
      The status should be failure
      The stderr should include 'required flag(s) "format" not set'
   End
End

Describe './snyk-iac-rules template ./fixtures/custom-rules --rule test --format fake'
   It 'returns  an error'
      When call ./snyk-iac-rules template ./fixtures/custom-rules --rule test --format fake
      The status should be failure
      The stderr should include 'invalid argument "fake" for "-f, --format" flag'
   End
End

Describe './snyk-iac-rules template ./fixtures/custom-rules --rule test1 --format hcl2'
   It 'generates the template'
      When call ./snyk-iac-rules template ./fixtures/custom-rules --rule test1 --format hcl2
      The status should be success
      The output should include 'Template rules directory'
      The output should include 'Template rules/test1 directory'
      The output should include 'Template rules/test1/main.rego file'
      The output should include 'Template rules/test1/main_test.rego file'
      The output should include 'Template rules/test1/fixtures directory'
      The output should include 'Template rules/test1/fixtures/denied.tf file'
      The output should include 'Template rules/test1/fixtures/allowed.tf file'
      The output should include 'Generated template'

      rm -rf ./fixtures/custom-rules/rules/test1
   End
End

Describe './snyk-iac-rules template ./fixtures/custom-rules --rule test2 --format json'
   It 'generates the template'
      When call ./snyk-iac-rules template ./fixtures/custom-rules --rule test2 --format json
      The status should be success
      The output should include 'Template rules directory'
      The output should include 'Template rules/test2 directory'
      The output should include 'Template rules/test2/main.rego file'
      The output should include 'Template rules/test2/main_test.rego file'
      The output should include 'Template rules/test2/fixtures directory'
      The output should include 'Template rules/test2/fixtures/denied.json file'
      The output should include 'Template rules/test2/fixtures/allowed.json file'
      The output should include 'Generated template'

      rm -rf ./fixtures/custom-rules/rules/test2
   End
End

Describe './snyk-iac-rules template ./fixtures/custom-rules --rule test3 --format yaml'
   It 'returns passing test status'
      When call ./snyk-iac-rules template ./fixtures/custom-rules --rule test3 --format yaml
      The status should be success
      The output should include 'Template rules directory'
      The output should include 'Template rules/test3 directory'
      The output should include 'Template rules/test3/main.rego file'
      The output should include 'Template rules/test3/main_test.rego file'
      The output should include 'Template rules/test3/fixtures directory'
      The output should include 'Template rules/test3/fixtures/denied.yaml file'
      The output should include 'Template rules/test3/fixtures/allowed.yaml file'
      The output should include 'Generated template'

      rm -rf ./fixtures/custom-rules/rules/test3
   End
End

Describe './snyk-iac-rules template ./fixtures/custom-rules --rule test4 --format tf-plan'
   It 'returns passing test status'
      When call ./snyk-iac-rules template ./fixtures/custom-rules --rule test4 --format tf-plan
      The status should be success
      The output should include 'Template rules directory'
      The output should include 'Template rules/test4 directory'
      The output should include 'Template rules/test4/main.rego file'
      The output should include 'Template rules/test4/main_test.rego file'
      The output should include 'Template rules/test4/fixtures directory'
      The output should include 'Template rules/test4/fixtures/denied.json.tfplan file'
      The output should include 'Template rules/test4/fixtures/allowed.json.tfplan file'
      The output should include 'Generated template'

      rm -rf ./fixtures/custom-rules/rules/test4
   End
End

Describe './snyk-iac-rules template ./fixtures/custom-rules --rule test4'
   It 'returns passing test status'
      ./snyk-iac-rules template ./fixtures/custom-rules --rule test4 --format tf-plan

      When call ./snyk-iac-rules template ./fixtures/custom-rules --rule test4 --format tf-plan
      The status should be failure
      The output should include 'Template rules directory'
      The stderr should include 'Rule with the provided name already exists'

      rm -rf ./fixtures/custom-rules/rules/test4
   End
End
