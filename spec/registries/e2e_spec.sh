#!/bin/bash

setup() {
    # create bundle
    ./snyk-iac-rules build ./fixtures/custom-rules --ignore "testing" --ignore "*_test.rego" 
}

cleanup() {
    rm bundle.tar.gz;
}

BeforeAll 'setup'
AfterAll 'cleanup'

Describe 'Supported Registry'
    registry_test() {
        OCI_REGISTRY_NAME="$1"
        OCI_REGISTRY_URL="$2"
        OCI_REGISTRY_USERNAME="$3"
        OCI_REGISTRY_PASSWORD="$4"
        ./snyk-iac-rules push --registry "$OCI_REGISTRY_NAME" bundle.tar.gz

        @registry_test "$OCI_REGISTRY_URL" "$OCI_REGISTRY_USERNAME" "$OCI_REGISTRY_PASSWORD" ./fixtures/custom-rules/rules/CUSTOM-1/fixtures/test.tf
    }
    Describe 'DockerHub'
        Skip if 'environment variable is not set'  [ -z "$OCI_DOCKERHUB_REGISTRY_URL" ]
        It "can push and pull"
            When call registry_test "$OCI_DOCKERHUB_REGISTRY_NAME" "$OCI_DOCKERHUB_REGISTRY_URL" "$OCI_DOCKERHUB_REGISTRY_USERNAME" "$OCI_DOCKERHUB_REGISTRY_PASSWORD"
            The status should eq 1
            The output should include 'Missing tags [Low Severity] [CUSTOM-1]' # it should include the custom rule in its output
        End
    End

    Describe "Azure"
        Skip if 'environment variable is not set'  [ -z "$OCI_AZURE_REGISTRY_URL" ]
        It "can push and pull"
            When call registry_test "$OCI_AZURE_REGISTRY_NAME" "$OCI_AZURE_REGISTRY_URL" "$OCI_AZURE_REGISTRY_USERNAME" "$OCI_AZURE_REGISTRY_PASSWORD"
            The status should eq 1
            The output should include 'Missing tags [Low Severity] [CUSTOM-1]' # it should include the custom rule in its output
        End
    End

End
