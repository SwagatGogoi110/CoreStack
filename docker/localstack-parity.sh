#!/bin/sh
# Maps LocalStack Community environment variables to their CloudStack equivalents.
# Sourced by entrypoint.sh when LOCALSTACK_PARITY=true.
# CloudStack vars always win: every mapping uses ${CLOUDSTACK_VAR:-<derived>} so an
# explicitly-set CloudStack var is never overwritten.

# Storage mode — PERSISTENCE=1 / PERSIST_STATE=1 → persistent storage
if [ -n "${PERSISTENCE:-}" ] || [ -n "${PERSIST_STATE:-}" ]; then
    _ls_persist="${PERSISTENCE:-${PERSIST_STATE:-}}"
    if [ "${_ls_persist}" = "1" ] || [ "${_ls_persist}" = "true" ]; then
        export CLOUDSTACK_STORAGE_MODE="${CLOUDSTACK_STORAGE_MODE:-persistent}"
    fi
fi

# Bind port — EDGE_PORT → CLOUDSTACK_PORT
[ -n "${EDGE_PORT:-}" ] && export CLOUDSTACK_PORT="${CLOUDSTACK_PORT:-${EDGE_PORT}}"

# Hostname returned in response URLs — LOCALSTACK_HOST / LOCALSTACK_HOSTNAME → CLOUDSTACK_HOSTNAME
_ls_host="${LOCALSTACK_HOST:-${LOCALSTACK_HOSTNAME:-}}"
[ -n "${_ls_host}" ] && export CLOUDSTACK_HOSTNAME="${CLOUDSTACK_HOSTNAME:-${_ls_host}}"

# Bind address — GATEWAY_LISTEN → CLOUDSTACK_HOST
[ -n "${GATEWAY_LISTEN:-}" ] && export CLOUDSTACK_HOST="${CLOUDSTACK_HOST:-${GATEWAY_LISTEN}}"

# Log level — LS_LOG / DEBUG=1 → LOG_LEVEL
if [ -n "${LS_LOG:-}" ]; then
    export LOG_LEVEL="${LOG_LEVEL:-${LS_LOG}}"
elif [ "${DEBUG:-}" = "1" ]; then
    export LOG_LEVEL="${LOG_LEVEL:-DEBUG}"
fi

# Lambda — LAMBDA_EXECUTOR is intentionally ignored; CloudStack always runs Lambda in Docker containers.

# Lambda Docker network
[ -n "${LAMBDA_DOCKER_NETWORK:-}" ] && \
    export CLOUDSTACK_SERVICES_LAMBDA_DOCKER_NETWORK="${CLOUDSTACK_SERVICES_LAMBDA_DOCKER_NETWORK:-${LAMBDA_DOCKER_NETWORK}}"

# Lambda ephemeral containers
if [ "${LAMBDA_REMOVE_CONTAINERS:-}" = "1" ] || [ "${LAMBDA_REMOVE_CONTAINERS:-}" = "true" ]; then
    export CLOUDSTACK_SERVICES_LAMBDA_EPHEMERAL="${CLOUDSTACK_SERVICES_LAMBDA_EPHEMERAL:-true}"
fi

# LAMBDA_REMOTE_DOCKER — not fully supported.
# CloudStack's hot-reload is per-function opt-in (S3Bucket=hot-reload), not a global bind-mount mode.
if [ -n "${LAMBDA_REMOTE_DOCKER:-}" ]; then
    echo "[CloudStack-parity] WARNING: LAMBDA_REMOTE_DOCKER is not fully supported by CloudStack." >&2
    echo "[CloudStack-parity] Use S3Bucket=hot-reload per function instead. See https://cloudstack.io/docs/services/lambda" >&2
fi

# Docker host
[ -n "${DOCKER_HOST:-}" ] && export CLOUDSTACK_DOCKER_DOCKER_HOST="${CLOUDSTACK_DOCKER_DOCKER_HOST:-${DOCKER_HOST}}"

# Docker network — shared across all container-based services (Lambda, RDS, ElastiCache, MSK, OpenSearch, EKS).
# Per-service overrides (e.g. CLOUDSTACK_SERVICES_LAMBDA_DOCKER_NETWORK) take precedence when set.
[ -n "${DOCKER_NETWORK:-}" ] && export CLOUDSTACK_SERVICES_DOCKER_NETWORK="${CLOUDSTACK_SERVICES_DOCKER_NETWORK:-${DOCKER_NETWORK}}"

# DNS suffixes — register LocalStack and CloudStack hostname suffixes so that container-to-container
# hostname routing (Function URLs, presigned S3, SQS QueueUrl, etc.) works without manual config.
_parity_suffixes="localhost.localstack.cloud,localhost.cloudstack.io"
if [ -n "${CLOUDSTACK_DNS_EXTRA_SUFFIXES:-}" ]; then
    export CLOUDSTACK_DNS_EXTRA_SUFFIXES="${CLOUDSTACK_DNS_EXTRA_SUFFIXES},${_parity_suffixes}"
else
    export CLOUDSTACK_DNS_EXTRA_SUFFIXES="${_parity_suffixes}"
fi

# TLS — USE_SSL=1 → CLOUDSTACK_TLS_ENABLED=true
if [ "${USE_SSL:-0}" = "1" ]; then
    export CLOUDSTACK_TLS_ENABLED="${CLOUDSTACK_TLS_ENABLED:-true}"
fi

# TLS — CUSTOM_SSL_CERT_PATH → CLOUDSTACK_TLS_CERT_PATH + CLOUDSTACK_TLS_KEY_PATH
# LocalStack uses a single combined PEM (cert+key). CloudStack accepts it in both fields.
if [ -n "${CUSTOM_SSL_CERT_PATH:-}" ]; then
    export CLOUDSTACK_TLS_CERT_PATH="${CLOUDSTACK_TLS_CERT_PATH:-${CUSTOM_SSL_CERT_PATH}}"
    export CLOUDSTACK_TLS_KEY_PATH="${CLOUDSTACK_TLS_KEY_PATH:-${CUSTOM_SSL_CERT_PATH}}"
fi

# SERVICES — intentionally ignored; CloudStack starts all 41 services in ~24ms.
