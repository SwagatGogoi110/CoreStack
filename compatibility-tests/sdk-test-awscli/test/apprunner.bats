#!/usr/bin/env bats
# AppRunner integration tests

setup() {
    load 'test_helper/common-setup'
    SVC_NAME="bats-ar-$(unique_name)"
}
teardown() {
    # Service ARN is needed for delete, we'd need to fetch it first
    true
}
@test "AppRunner: create and describe service" {
    run aws_cmd apprunner create-service --service-name "$SVC_NAME"
    assert_success
    SVC_ARN=$(echo "$output" | jq -r '.Service.ServiceArn')
    run aws_cmd apprunner describe-service --service-arn "$SVC_ARN"
    assert_success
    [ "$(echo "$output" | jq -r '.Service.ServiceName')" = "$SVC_NAME" ]
}