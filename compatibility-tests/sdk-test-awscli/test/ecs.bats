#!/usr/bin/env bats
# ECS integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "ECS: basic list/describe operation" {
    run aws_cmd ecs list-clusters
    assert_success
}
