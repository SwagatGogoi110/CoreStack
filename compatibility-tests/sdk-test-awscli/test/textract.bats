#!/usr/bin/env bats
# Textract integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "Textract: basic list/describe operation" {
    run aws_cmd textract list-adapters
    assert_success
}
