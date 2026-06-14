#!/usr/bin/env bats
# Firehose integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "Firehose: basic list/describe operation" {
    run aws_cmd firehose list-delivery-streams
    assert_success
}
