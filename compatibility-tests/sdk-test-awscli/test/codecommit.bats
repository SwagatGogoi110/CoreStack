#!/usr/bin/env bats
# CodeCommit integration tests

setup() {
    load 'test_helper/common-setup'
    REPO_NAME="bats-repo-$(unique_name)"
}
teardown() {
    aws_cmd codecommit delete-repository --repository-name "$REPO_NAME" >/dev/null 2>&1 || true
}
@test "CodeCommit: create and get repository" {
    run aws_cmd codecommit create-repository --repository-name "$REPO_NAME"
    assert_success
    run aws_cmd codecommit get-repository --repository-name "$REPO_NAME"
    assert_success
    [ "$(echo "$output" | jq -r '.repositoryMetadata.repositoryName')" = "$REPO_NAME" ]
}