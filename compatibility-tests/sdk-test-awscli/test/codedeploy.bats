#!/usr/bin/env bats
# CodeDeploy integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "CodeDeploy: basic list/describe operation" {
    run aws_cmd codedeploy list-applications
    assert_success
}
