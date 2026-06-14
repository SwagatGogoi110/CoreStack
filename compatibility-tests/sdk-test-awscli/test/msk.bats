#!/usr/bin/env bats
# MSK integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "MSK: basic list/describe operation" {
    run aws_cmd kafka list-clusters
    assert_success
}
