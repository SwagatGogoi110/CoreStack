#!/usr/bin/env bats
# Transfer integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "Transfer: basic list/describe operation" {
    run aws_cmd transfer list-servers
    assert_success
}
