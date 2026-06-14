#!/usr/bin/env bats
# API Gateway integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "API Gateway: basic list/describe operation" {
    run aws_cmd apigateway get-rest-apis
    assert_success
}
