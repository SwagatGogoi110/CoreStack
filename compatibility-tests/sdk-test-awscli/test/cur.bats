#!/usr/bin/env bats
# Cost and Usage Report integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "Cost and Usage Report: basic list/describe operation" {
    run aws_cmd cur describe-report-definitions
    assert_success
}
