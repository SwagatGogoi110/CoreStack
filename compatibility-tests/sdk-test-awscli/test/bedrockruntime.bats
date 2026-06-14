#!/usr/bin/env bats
# Bedrock Runtime integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "Bedrock Runtime: basic list/describe operation" {
    run aws_cmd bedrock-runtime list-custom-models
    assert_success
}
