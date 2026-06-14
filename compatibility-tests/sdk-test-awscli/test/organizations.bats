#!/usr/bin/env bats
# Organizations integration tests

setup() {
    load 'test_helper/common-setup'
}
@test "Organizations: create and describe organization" {
    run aws_cmd organizations create-organization --feature-set ALL
    # Might already exist if run multiple times without reset
    run aws_cmd organizations describe-organization
    assert_success
    [ "$(echo "$output" | jq -r '.Organization.FeatureSet')" = "ALL" ]
}