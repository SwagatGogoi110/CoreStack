#!/usr/bin/env bats
# EMR integration tests

setup() {
    load 'test_helper/common-setup'
    CLUSTER_NAME="bats-emr-$(unique_name)"
}
@test "EMR: run job flow and list clusters" {
    run aws_cmd emr create-cluster --name "$CLUSTER_NAME" --release-label emr-6.2.0 --use-default-roles --instance-type m5.xlarge --instance-count 1
    assert_success
    JOB_ID=$(echo "$output" | jq -r '.ClusterId')
    run aws_cmd emr describe-cluster --cluster-id "$JOB_ID"
    assert_success
    [ "$(echo "$output" | jq -r '.Cluster.Name')" = "$CLUSTER_NAME" ]
}