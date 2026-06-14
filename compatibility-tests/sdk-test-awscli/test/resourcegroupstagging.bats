#!/usr/bin/env bats
# Resource Groups Tagging integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "Resource Groups Tagging: basic list/describe operation" {
    run aws_cmd resourcegroupstaggingapi get-resources
    assert_success
}
