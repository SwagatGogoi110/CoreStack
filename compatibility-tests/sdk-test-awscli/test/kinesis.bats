#!/usr/bin/env bats
# Kinesis integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "Kinesis: basic list/describe operation" {
    run aws_cmd kinesis list-streams
    assert_success
}
