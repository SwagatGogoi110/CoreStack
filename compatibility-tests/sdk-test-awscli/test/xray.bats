#!/usr/bin/env bats
# X-Ray integration tests

setup() {
    load 'test_helper/common-setup'
}
@test "X-Ray: put segments and get summaries" {
    run aws_cmd xray put-trace-segments --trace-segment-documents "{\"id\": \"1\", \"name\": \"myservice\"}"
    assert_success
    run aws_cmd xray get-trace-summaries --start-time 1 --end-time 2
    assert_success
}