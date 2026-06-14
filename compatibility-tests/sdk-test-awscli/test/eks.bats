#!/usr/bin/env bats
# EKS integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "EKS: basic list/describe operation" {
    run aws_cmd eks list-clusters
    assert_success
}
