#!/bin/bash

cleanupBundle() { rm bundle.tar.gz; }
AfterAll 'cleanupBundle'

setupTmp() { mkdir tmp; }
cleanupTmp() { rm -rf ./tmp; }

Describe 'Contract test between the SDK and the Snyk CLI'
    Describe 'Via --rules flag'
        BeforeEach setupTmp
        AfterEach cleanupTmp
        It 'Verifies custom rule without a path'
            snyk_iac_test() {
                cd tmp

                # create a basic rule
                ../snyk-iac-rules template --rule Contract --format hcl2

                # run tests and make sure they pass
                ../snyk-iac-rules test

                # create bundle
                ../snyk-iac-rules build --ignore "testing" --ignore "*_test.rego" 

                # authenticate with Snyk
                snyk auth $SNYK_TOKEN 

                # use bundle with Snyk
                snyk iac test --rules=./bundle.tar.gz ./rules/Contract/fixtures/denied.tf
                echo $?
            }

            When call snyk_iac_test
            The status should eq 0
            The output should include "Generated template" # the rule was tempalted successfully
            The output should include "PASS: 1/1" # the tests passed
            The output should include "Generated bundle: bundle.tar.gz" # the bundle has been generated
            The output should include "Using custom rules to generate misconfigurations." # uses the custom rules to generate misconfigurations
            The output should include "Default title [Low Severity] [Contract]" # it should include the custom rule in its output
            The output should not include 'WARNING: The command must point at a folder that contains the package for the rules'
            The stderr should not be present

            cd ../
        End

        It 'Verifies custom rule with relative path'
            snyk_iac_test() {
                # create a basic rule
                ./snyk-iac-rules template ./tmp --rule Contract --format hcl2

                OS=$(uname)
                # replace the fixture path so it's correct
                if [ "$OS" = "Darwin" ]; then
                    sed -i '' -e 's#/rules/Contract/fixtures#/tmp/rules/Contract/fixtures#' ./tmp/rules/Contract/main_test.rego
                else
                    sed -i -e 's#/rules/Contract/fixtures#/tmp/rules/Contract/fixtures#' ./tmp/rules/Contract/main_test.rego
                fi

                # run tests and make sure they pass
                ./snyk-iac-rules test ./tmp 

                # create bundle
                ./snyk-iac-rules build ./tmp --ignore "testing" --ignore "*_test.rego"

                # authenticate with Snyk
                snyk auth $SNYK_TOKEN

                # use bundle with Snyk
                snyk iac test --rules=./bundle.tar.gz ./tmp/rules/Contract/fixtures/denied.tf
                echo $?
            }

            When call snyk_iac_test
            The status should eq 0
            The output should include "Generated template" # the rule was tempalted successfully
            The output should include "PASS: 1/1" # the tests passed
            The output should include "Generated bundle: bundle.tar.gz" # the bundle has been generated
            The output should include "Using custom rules to generate misconfigurations." # uses the custom rules to generate misconfigurations
            The output should include "Default title [Low Severity] [Contract]" # it should include the custom rule in its output
            The output should not include 'WARNING: The command must point at a folder that contains the package for the rules'
            The stderr should not be present
        End
    End

    Describe 'Via push and pull'
        skip_push_test() { ! [ -z "$SKIP_PUSH_TEST" ]; }
        Skip if 'skip environment variable is set' skip_push_test
        BeforeEach setupTmp
        AfterEach cleanupTmp
        It 'verifies contract between the SDK and Snyk'
            snyk_iac_test() {
                cd tmp

                # create a basic rule
                ../snyk-iac-rules template --rule Contract --format hcl2

                # run tests and make sure they pass
                ../snyk-iac-rules test

                # create bundle
                ../snyk-iac-rules build --ignore "testing" --ignore "*_test.rego" 

                # push bundle
                ../snyk-iac-rules push --registry $OCI_REGISTRY_NAME-$OS bundle.tar.gz

                @registry_test https://registry-1.$OCI_REGISTRY_NAME-$OS $OCI_REGISTRY_USERNAME $OCI_REGISTRY_PASSWORD ./rules/Contract/fixtures/denied.tf
                echo $?
            }

            When call snyk_iac_test
            The status should eq 0
            The output should include "Generated template" # the rule was tempalted successfully
            The output should include "PASS: 1/1" # the tests passed
            The output should include "Generated bundle: bundle.tar.gz" # the bundle has been generated
            The output should include "Using custom rules to generate misconfigurations." # uses the custom rules to generate misconfigurations
            The output should include "Default title [Low Severity] [Contract]" # it should include the custom rule in its output
            The output should not include 'WARNING: The command must point at a folder that contains the package for the rules'
            The stderr should not be present

            cd ../
        End
    End
End
