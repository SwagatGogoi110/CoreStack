#!/usr/bin/env bats
# Cost Explorer integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "Cost Explorer: basic list/describe operation" {
    run aws_cmd ce get-cost-and-usage
    assert_success
}
