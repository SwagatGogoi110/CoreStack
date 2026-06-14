#!/usr/bin/env bats
# CloudFront integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "CloudFront: basic list/describe operation" {
    run aws_cmd cloudfront list-distributions
    assert_success
}
