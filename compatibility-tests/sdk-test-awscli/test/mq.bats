#!/usr/bin/env bats
# Amazon MQ integration tests

setup() {
    load 'test_helper/common-setup'
    BROKER_NAME="bats-mq-$(unique_name)"
}
@test "MQ: create and describe broker" {
    run aws_cmd mq create-broker --broker-name "$BROKER_NAME" --engine-type ACTIVEMQ --deployment-mode SINGLE_INSTANCE
    assert_success
    BROKER_ID=$(echo "$output" | jq -r '.brokerId')
    run aws_cmd mq describe-broker --broker-id "$BROKER_ID"
    assert_success
    [ "$(echo "$output" | jq -r '.brokerName')" = "$BROKER_NAME" ]
}