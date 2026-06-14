#!/usr/bin/env bats
# Glue integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "Glue: basic list/describe operation" {
    run aws_cmd glue get-databases
    assert_success
}
