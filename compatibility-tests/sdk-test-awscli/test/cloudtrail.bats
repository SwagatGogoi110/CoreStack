#!/usr/bin/env bats
# CloudTrail integration tests

setup() {
    load 'test_helper/common-setup'
    TRAIL_NAME="bats-trail-$(unique_name)"
}
teardown() {
    aws_cmd cloudtrail delete-trail --name "$TRAIL_NAME" >/dev/null 2>&1 || true
}
@test "CloudTrail: create and describe trail" {
    run aws_cmd cloudtrail create-trail --name "$TRAIL_NAME" --s3-bucket-name mybucket
    assert_success
    run aws_cmd cloudtrail describe-trails --trail-name-list "$TRAIL_NAME"
    assert_success
    [ "$(echo "$output" | jq -r '.trailList[0].Name')" = "$TRAIL_NAME" ]
}