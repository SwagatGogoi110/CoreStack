#!/usr/bin/env bats
# Athena integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "Athena: basic list/describe operation" {
    run aws_cmd athena list-work-groups
    assert_success
}
