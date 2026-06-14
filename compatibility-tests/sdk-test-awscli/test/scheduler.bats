#!/usr/bin/env bats
# Scheduler integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "Scheduler: basic list/describe operation" {
    run aws_cmd scheduler list-schedules
    assert_success
}
