#!/usr/bin/env bats
# EventBridge integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "EventBridge: basic list/describe operation" {
    run aws_cmd events list-event-buses
    assert_success
}
