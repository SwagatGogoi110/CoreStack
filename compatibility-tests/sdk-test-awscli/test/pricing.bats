#!/usr/bin/env bats
# Pricing integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "Pricing: basic list/describe operation" {
    run aws_cmd pricing describe-services
    assert_success
}
