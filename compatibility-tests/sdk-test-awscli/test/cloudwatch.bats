#!/usr/bin/env bats
# CloudWatch integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "CloudWatch: basic list/describe operation" {
    run aws_cmd cloudwatch list-metrics
    assert_success
}
