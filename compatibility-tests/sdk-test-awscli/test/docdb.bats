#!/usr/bin/env bats
# DocumentDB integration tests

setup() {
    load 'test_helper/common-setup'
    CLUSTER_ID="bats-docdb-$(unique_name)"
}
teardown() {
    aws_cmd docdb delete-db-cluster --db-cluster-identifier "$CLUSTER_ID" >/dev/null 2>&1 || true
}
@test "DocumentDB: create and describe cluster" {
    run aws_cmd docdb create-db-cluster --db-cluster-identifier "$CLUSTER_ID" --engine docdb
    assert_success
    run aws_cmd docdb describe-db-clusters --db-cluster-identifier "$CLUSTER_ID"
    assert_success
    [ "$(echo "$output" | jq -r '.DBClusters[0].DBClusterIdentifier')" = "$CLUSTER_ID" ]
}