#!/usr/bin/env bats
# Transcribe integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "Transcribe: basic list/describe operation" {
    run aws_cmd transcribe list-transcription-jobs
    assert_success
}
