#!/usr/bin/env bats
# ELBv2 integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "ELBv2: basic list/describe operation" {
    run aws_cmd elbv2 describe-load-balancers
    assert_success
}
