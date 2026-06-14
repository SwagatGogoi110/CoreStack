#!/usr/bin/env bats
# Redshift integration tests

setup() {
    load 'test_helper/common-setup'
    CLUSTER_ID="bats-rs-$(unique_name)"
}
teardown() {
    aws_cmd redshift delete-cluster --cluster-identifier "$CLUSTER_ID" >/dev/null 2>&1 || true
}
@test "Redshift: create and describe cluster" {
    run aws_cmd redshift create-cluster --cluster-identifier "$CLUSTER_ID" --node-type dc2.large --master-username admin
    assert_success
    run aws_cmd redshift describe-clusters --cluster-identifier "$CLUSTER_ID"
    assert_success
    [ "$(echo "$output" | jq -r '.Clusters[0].ClusterIdentifier')" = "$CLUSTER_ID" ]
}