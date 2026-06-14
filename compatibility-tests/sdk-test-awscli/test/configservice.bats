#!/usr/bin/env bats
# Config Service integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "Config Service: basic list/describe operation" {
    run aws_cmd configservice describe-config-rules
    assert_success
}
