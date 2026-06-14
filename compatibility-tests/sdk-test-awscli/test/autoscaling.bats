#!/usr/bin/env bats
# Auto Scaling integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "Auto Scaling: basic list/describe operation" {
    run aws_cmd autoscaling describe-auto-scaling-groups
    assert_success
}
