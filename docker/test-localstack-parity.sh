#!/bin/sh
# Unit tests for localstack-parity.sh.
# Run directly: sh docker/test-localstack-parity.sh
# Exit 0 on success, non-zero on first failure.

set -eu

SCRIPT="$(dirname "$0")/localstack-parity.sh"
PASS=0
FAIL=0

# Run the parity script in a subshell with a given environment and print the
# value of a single variable. Arguments: VAR_NAME [ENV_KEY=VALUE ...]
_run() {
    var="$1"; shift
    env -i "$@" sh -c ". '${SCRIPT}'; printf '%s' \"\${${var}:-}\""
}

# Assert that _run produces an expected value.
assert_eq() {
    desc="$1"; expected="$2"; actual="$3"
    if [ "${actual}" = "${expected}" ]; then
        printf '[PASS] %s\n' "${desc}"
        PASS=$((PASS + 1))
    else
        printf '[FAIL] %s\n  expected: %s\n  actual:   %s\n' "${desc}" "${expected}" "${actual}"
        FAIL=$((FAIL + 1))
    fi
}

# --- PERSISTENCE ---
assert_eq "PERSISTENCE=1 sets CLOUDSTACK_STORAGE_MODE=persistent" \
    "persistent" \
    "$(_run CLOUDSTACK_STORAGE_MODE PERSISTENCE=1)"

assert_eq "PERSISTENCE=true sets CLOUDSTACK_STORAGE_MODE=persistent" \
    "persistent" \
    "$(_run CLOUDSTACK_STORAGE_MODE PERSISTENCE=true)"

assert_eq "PERSIST_STATE=1 sets CLOUDSTACK_STORAGE_MODE=persistent" \
    "persistent" \
    "$(_run CLOUDSTACK_STORAGE_MODE PERSIST_STATE=1)"

assert_eq "CLOUDSTACK_STORAGE_MODE wins over PERSISTENCE" \
    "hybrid" \
    "$(_run CLOUDSTACK_STORAGE_MODE PERSISTENCE=1 CLOUDSTACK_STORAGE_MODE=hybrid)"

# --- EDGE_PORT ---
assert_eq "EDGE_PORT sets CLOUDSTACK_PORT" \
    "4567" \
    "$(_run CLOUDSTACK_PORT EDGE_PORT=4567)"

assert_eq "CLOUDSTACK_PORT wins over EDGE_PORT" \
    "4568" \
    "$(_run CLOUDSTACK_PORT EDGE_PORT=4567 CLOUDSTACK_PORT=4568)"

# --- LOCALSTACK_HOST / LOCALSTACK_HOSTNAME ---
assert_eq "LOCALSTACK_HOST sets CLOUDSTACK_HOSTNAME" \
    "myhost" \
    "$(_run CLOUDSTACK_HOSTNAME LOCALSTACK_HOST=myhost)"

assert_eq "LOCALSTACK_HOSTNAME sets CLOUDSTACK_HOSTNAME when LOCALSTACK_HOST unset" \
    "myhost2" \
    "$(_run CLOUDSTACK_HOSTNAME LOCALSTACK_HOSTNAME=myhost2)"

assert_eq "LOCALSTACK_HOST takes priority over LOCALSTACK_HOSTNAME" \
    "primary" \
    "$(_run CLOUDSTACK_HOSTNAME LOCALSTACK_HOST=primary LOCALSTACK_HOSTNAME=secondary)"

assert_eq "CLOUDSTACK_HOSTNAME wins over LOCALSTACK_HOST" \
    "explicit" \
    "$(_run CLOUDSTACK_HOSTNAME LOCALSTACK_HOST=myhost CLOUDSTACK_HOSTNAME=explicit)"

# --- GATEWAY_LISTEN ---
assert_eq "GATEWAY_LISTEN sets QUARKUS_HTTP_HOST" \
    "0.0.0.0" \
    "$(_run QUARKUS_HTTP_HOST GATEWAY_LISTEN=0.0.0.0)"

# --- LOG LEVEL ---
assert_eq "LS_LOG sets QUARKUS_LOG_LEVEL" \
    "WARN" \
    "$(_run QUARKUS_LOG_LEVEL LS_LOG=WARN)"

assert_eq "DEBUG=1 sets QUARKUS_LOG_LEVEL=DEBUG" \
    "DEBUG" \
    "$(_run QUARKUS_LOG_LEVEL DEBUG=1)"

assert_eq "LS_LOG takes priority over DEBUG=1" \
    "TRACE" \
    "$(_run QUARKUS_LOG_LEVEL LS_LOG=TRACE DEBUG=1)"

assert_eq "QUARKUS_LOG_LEVEL wins over LS_LOG" \
    "INFO" \
    "$(_run QUARKUS_LOG_LEVEL LS_LOG=DEBUG QUARKUS_LOG_LEVEL=INFO)"

# --- LAMBDA ---
assert_eq "LAMBDA_DOCKER_NETWORK sets CLOUDSTACK_SERVICES_LAMBDA_DOCKER_NETWORK" \
    "mynet" \
    "$(_run CLOUDSTACK_SERVICES_LAMBDA_DOCKER_NETWORK LAMBDA_DOCKER_NETWORK=mynet)"

assert_eq "CLOUDSTACK_SERVICES_LAMBDA_DOCKER_NETWORK wins over LAMBDA_DOCKER_NETWORK" \
    "CloudStack-net" \
    "$(_run CLOUDSTACK_SERVICES_LAMBDA_DOCKER_NETWORK LAMBDA_DOCKER_NETWORK=mynet CLOUDSTACK_SERVICES_LAMBDA_DOCKER_NETWORK=CloudStack-net)"

assert_eq "LAMBDA_REMOVE_CONTAINERS=1 sets CLOUDSTACK_SERVICES_LAMBDA_EPHEMERAL=true" \
    "true" \
    "$(_run CLOUDSTACK_SERVICES_LAMBDA_EPHEMERAL LAMBDA_REMOVE_CONTAINERS=1)"

assert_eq "LAMBDA_REMOVE_CONTAINERS=true sets CLOUDSTACK_SERVICES_LAMBDA_EPHEMERAL=true" \
    "true" \
    "$(_run CLOUDSTACK_SERVICES_LAMBDA_EPHEMERAL LAMBDA_REMOVE_CONTAINERS=true)"

assert_eq "CLOUDSTACK_SERVICES_LAMBDA_EPHEMERAL wins over LAMBDA_REMOVE_CONTAINERS" \
    "false" \
    "$(_run CLOUDSTACK_SERVICES_LAMBDA_EPHEMERAL LAMBDA_REMOVE_CONTAINERS=1 CLOUDSTACK_SERVICES_LAMBDA_EPHEMERAL=false)"

# --- DOCKER HOST / NETWORK ---
assert_eq "DOCKER_HOST sets CLOUDSTACK_DOCKER_DOCKER_HOST" \
    "unix:///var/run/docker.sock" \
    "$(_run CLOUDSTACK_DOCKER_DOCKER_HOST DOCKER_HOST=unix:///var/run/docker.sock)"

assert_eq "DOCKER_NETWORK sets CLOUDSTACK_SERVICES_DOCKER_NETWORK" \
    "shared" \
    "$(_run CLOUDSTACK_SERVICES_DOCKER_NETWORK DOCKER_NETWORK=shared)"

assert_eq "CLOUDSTACK_SERVICES_DOCKER_NETWORK wins over DOCKER_NETWORK" \
    "override" \
    "$(_run CLOUDSTACK_SERVICES_DOCKER_NETWORK DOCKER_NETWORK=shared CLOUDSTACK_SERVICES_DOCKER_NETWORK=override)"

# --- DNS SUFFIXES ---
assert_eq "DNS suffixes set when CLOUDSTACK_DNS_EXTRA_SUFFIXES unset" \
    "localhost.localstack.cloud,localhost.cloudstack.io" \
    "$(_run CLOUDSTACK_DNS_EXTRA_SUFFIXES)"

assert_eq "DNS suffixes appended to existing CLOUDSTACK_DNS_EXTRA_SUFFIXES" \
    "custom.internal,localhost.localstack.cloud,localhost.cloudstack.io" \
    "$(_run CLOUDSTACK_DNS_EXTRA_SUFFIXES CLOUDSTACK_DNS_EXTRA_SUFFIXES=custom.internal)"

# --- Summary ---
printf '\nResults: %d passed, %d failed\n' "${PASS}" "${FAIL}"
[ "${FAIL}" -eq 0 ]
