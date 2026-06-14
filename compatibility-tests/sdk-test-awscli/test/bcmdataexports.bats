#!/usr/bin/env bats
# BCM Data Exports integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "BCM Data Exports: basic list/describe operation" {
    run aws_cmd bcmdataexports list-exports
    assert_success
}
