#!/usr/bin/env bats
# Route53 integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "Route53: basic list/describe operation" {
    run aws_cmd route53 list-hosted-zones
    assert_success
}
