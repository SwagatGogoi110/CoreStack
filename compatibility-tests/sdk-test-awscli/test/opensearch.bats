#!/usr/bin/env bats
# OpenSearch integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "OpenSearch: basic list/describe operation" {
    run aws_cmd opensearch list-domain-names
    assert_success
}
