#!/usr/bin/env bats
# AppSync integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "AppSync: basic list/describe operation" {
    run aws_cmd appsync list-graphql-apis
    assert_success
}
