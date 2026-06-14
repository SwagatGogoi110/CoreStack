#!/usr/bin/env bats
# API Gateway V2 integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "API Gateway V2: basic list/describe operation" {
    run aws_cmd apigatewayv2 get-apis
    assert_success
}
