#!/usr/bin/env bats
# Synthetics integration tests

setup() {
    load 'test_helper/common-setup'
    CANARY_NAME="bats-canary-$(unique_name)"
}
teardown() {
    aws_cmd synthetics delete-canary --name "$CANARY_NAME" >/dev/null 2>&1 || true
}
@test "Synthetics: create and get canary" {
    run aws_cmd synthetics create-canary --name "$CANARY_NAME" --artifact-s3-location s3://mybucket --execution-role-arn arn:aws:iam::000000000000:role/canary
    assert_success
    run aws_cmd synthetics get-canary --name "$CANARY_NAME"
    assert_success
    [ "$(echo "$output" | jq -r '.Canary.Name')" = "$CANARY_NAME" ]
}