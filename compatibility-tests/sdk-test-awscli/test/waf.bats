#!/usr/bin/env bats
# WAF integration tests

setup() {
    load 'test_helper/common-setup'
    ACL_NAME="bats-waf-$(unique_name)"
}
@test "WAF: create and list web acls" {
    run aws_cmd waf create-web-acl --name "$ACL_NAME" --metric-name "$ACL_NAME" --default-action Type=ALLOW
    assert_success
    ACL_ID=$(echo "$output" | jq -r '.WebACL.WebACLId')
    run aws_cmd waf get-web-acl --web-acl-id "$ACL_ID"
    assert_success
    [ "$(echo "$output" | jq -r '.WebACL.Name')" = "$ACL_NAME" ]
}