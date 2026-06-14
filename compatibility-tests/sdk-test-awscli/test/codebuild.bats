#!/usr/bin/env bats
# CodeBuild integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "CodeBuild: basic list/describe operation" {
    run aws_cmd codebuild list-projects
    assert_success
}
