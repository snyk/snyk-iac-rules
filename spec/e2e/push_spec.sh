#!/bin/bash

setup() {
   ./snyk-iac-rules build ./fixtures/custom-rules --ignore testing --ignore "*_test.rego"
}
cleanup() { rm bundle.tar.gz; }

BeforeEach 'setup'
AfterEach 'cleanup'

Describe './snyk-iac-rules push -r docker.io/test/test test.jpg'
   It 'returns failing test status'
      When call ./snyk-iac-rules push -r docker.io/test/test test.jpg
      The status should be failure
      The output should include 'The path must be to a generated .tar.gz bundle'
      The stderr should be present
   End
End

# This tries to push a non-existant bundle to a DockerHub container registry
Describe './snyk-iac-rules push -r docker.io/test/test bundle-incorrect.tar.gz'
   It 'returns failing test status'
      When call ./snyk-iac-rules push -r docker.io/test/test bundle-incorrect.tar.gz
      The status should be failure
      The output should include 'Failed to read from the provided path'
      The stderr should be present
   End
End

Describe './snyk-iac-rules push -r test bundle.tar.gz'
   It 'returns failing test status'
      When call ./snyk-iac-rules push -r test bundle.tar.gz
      The status should be failure
      The output should include 'The provided container registry is invalid'
      The stderr should be present
   End
End

# This actually tries to push to a DockerHub container registry but we don't authenticate so it fails
Describe './snyk-iac-rules push -r docker.io/test/test bundle.tar.gz'
   It 'returns failing test status'
      When call ./snyk-iac-rules push -r docker.io/test/test bundle.tar.gz
      The status should be failure
      The output should include 'bundle.tar.gz'
      The output should include 'config.json'
      The output should include 'Failed to push bundle to container registry: server message: insufficient_scope: authorization failed'
      The stderr should be present
   End
End

Describe './snyk-iac-rules push -r $OCI_REGISTRY_NAME bundle.tar.gz'
   It 'returns passing test status'
      skip_push_test() { ! [ -z "$SKIP_PUSH_TEST" ]; }
      Skip if 'skip environment variable is set' skip_push_test
      When call ./snyk-iac-rules push -r $OCI_REGISTRY_NAME bundle.tar.gz
      The status should be success
      The output should include 'bundle.tar.gz'
      The output should include 'config.json'
      The output should include 'Successfully pushed bundle to'
   End
End

