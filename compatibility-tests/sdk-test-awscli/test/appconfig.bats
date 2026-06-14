#!/usr/bin/env bats
# AppConfig integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "AppConfig: basic list/describe operation" {
    run aws_cmd appconfig list-applications
    assert_success
}
