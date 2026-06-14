#!/usr/bin/env bats
# ElastiCache integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "ElastiCache: basic list/describe operation" {
    run aws_cmd elasticache describe-cache-clusters
    assert_success
}
