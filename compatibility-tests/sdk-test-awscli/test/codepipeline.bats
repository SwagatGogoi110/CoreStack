#!/usr/bin/env bats
# CodePipeline integration tests

setup() {
    load 'test_helper/common-setup'
    PIPE_NAME="bats-pipe-$(unique_name)"
}
teardown() {
    aws_cmd codepipeline delete-pipeline --name "$PIPE_NAME" >/dev/null 2>&1 || true
}
@test "CodePipeline: create and get pipeline" {
    run aws_cmd codepipeline create-pipeline --pipeline "{\"name\": \"$PIPE_NAME\", \"roleArn\": \"arn:aws:iam::000000000000:role/pipe\", \"stages\": []}"
    assert_success
    run aws_cmd codepipeline get-pipeline --name "$PIPE_NAME"
    assert_success
    [ "$(echo "$output" | jq -r '.pipeline.name')" = "$PIPE_NAME" ]
}