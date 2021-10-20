#!/bin/bash

cleanupBundle() { rm bundle.tar.gz; }
AfterAll 'cleanupBundle'

cleanup() { rm -rf ./tmp; unset OCI_REGISTRY_URL; }
AfterEach 'cleanup'

Describe 'Contract test'
    It 'verifies contract between the SDK and Snyk via flag'
        snyk_iac_test() {
            # create tmp test directory for contract tests
            mkdir tmp
            
            # create a basic rule
            ./snyk-iac-rules template ./tmp --rule Contract

            # replace the fixture path so it's correct
            sed -i '' -e 's#/rules/Contract/fixtures#/tmp/rules/Contract/fixtures#' ./tmp/rules/Contract/main_test.rego

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
        The output should include "Default title [Low Severity] [Contract]" # it should include the custom rule in its output
        The stderr should be present # from the templating
    End

    It 'verifies contract between the SDK and Snyk via distribution'
        skip_push_test() { ! [ -z "$SKIP_PUSH_TEST" ]; }
        Skip if 'skip environment variable is set' skip_push_test
        snyk_iac_test() {
            # create tmp test directory for contract tests
            mkdir tmp
            
            # create a basic rule
            ./snyk-iac-rules template ./tmp --rule Contract

            # replace the fixture path so it's correct
            sed -i '' -e 's#/rules/Contract/fixtures#/tmp/rules/Contract/fixtures#' ./tmp/rules/Contract/main_test.rego

            # run tests and make sure they pass
            ./snyk-iac-rules test ./tmp 

            # create bundle
            ./snyk-iac-rules build ./tmp --ignore "testing" --ignore "*_test.rego" 

            # push bundle
            ./snyk-iac-rules push --registry $OCI_REGISTRY_NAME bundle.tar.gz

            # authenticate with Snyk
            snyk auth $SNYK_TOKEN 

            # set environment variables for the CLI
            export OCI_REGISTRY_URL=https://registry-1.$OCI_REGISTRY_NAME
            export OCI_REGISTRY_USERNAME=$OCI_REGISTRY_USERNAME
            export OCI_REGISTRY_PASSWORD=$OCI_REGISTRY_PASSWORD

            # use bundle with Snyk
            snyk iac test ./tmp/rules/Contract/fixtures/denied.tf
            echo $?
        }

        When call snyk_iac_test
        The status should eq 0
        The output should include "Generated template" # the rule was tempalted successfully
        The output should include "PASS: 1/1" # the tests passed
        The output should include "Generated bundle: bundle.tar.gz" # the bundle has been generated
        The output should include "Default title [Low Severity] [Contract]" # it should include the custom rule in its output
        The stderr should be present # from the templating
    End
End
