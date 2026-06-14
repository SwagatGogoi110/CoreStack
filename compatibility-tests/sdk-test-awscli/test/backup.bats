#!/usr/bin/env bats
# Backup integration tests

setup() {
    load 'test_helper/common-setup'
}

@test "Backup: basic list/describe operation" {
    run aws_cmd backup list-backup-plans
    assert_success
}
