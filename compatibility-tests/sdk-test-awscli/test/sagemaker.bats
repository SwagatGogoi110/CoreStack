#!/usr/bin/env bats
# SageMaker integration tests

setup() {
    load 'test_helper/common-setup'
    NB_NAME="bats-nb-$(unique_name)"
}
@test "SageMaker: create and list notebook instances" {
    run aws_cmd sagemaker create-notebook-instance --notebook-instance-name "$NB_NAME" --instance-type ml.t2.medium
    assert_success
    run aws_cmd sagemaker list-notebook-instances
    assert_success
    [ "$(echo "$output" | jq -r '.NotebookInstances[0].NotebookInstanceName')" = "$NB_NAME" ]
}