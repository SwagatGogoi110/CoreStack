#!/usr/bin/env bats
# EC2 integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "EC2: basic list/describe operation" {
    run aws_cmd ec2 describe-instances
    assert_success
}
