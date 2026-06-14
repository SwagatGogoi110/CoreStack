#!/usr/bin/env bats
# Step Functions integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "Step Functions: basic list/describe operation" {
    run aws_cmd stepfunctions list-state-machines
    assert_success
}
