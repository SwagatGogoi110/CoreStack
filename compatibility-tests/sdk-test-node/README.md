# sdk-test-node

Compatibility tests for [CloudStack](https://github.com/hectorvent/CloudStack) using the **AWS SDK for JavaScript v3 (3.1003.0)**.

## Services Covered

| Group                   | Description                                                                                        |
| ----------------------- | -------------------------------------------------------------------------------------------------- |
| `ssm`                   | Parameter Store — put, get, label, history, path, tags                                             |
| `sqs`                   | Queues, send/receive/delete, DLQ, visibility                                                       |
| `sns`                   | Topics, subscriptions, publish, SQS delivery                                                       |
| `s3`                    | Buckets, objects, tagging, copy, batch delete                                                      |
| `s3-cors`               | CORS configuration                                                                                 |
| `s3-notifications`      | S3 → SQS and S3 → SNS event notifications                                                          |
| `dynamodb`              | Tables, CRUD, batch, TTL, tags                                                                     |
| `lambda`                | Create/invoke/update/delete functions                                                              |
| `iam`                   | Users, roles, policies, access keys                                                                |
| `sts`                   | GetCallerIdentity, AssumeRole, GetSessionToken                                                     |
| `secretsmanager`        | Create/get/put/list/delete secrets, versioning, tags                                               |
| `kms`                   | Keys, aliases, encrypt/decrypt, data keys, sign/verify                                             |
| `kinesis`               | Streams, shards, PutRecord/GetRecords                                                              |
| `cloudwatch`            | PutMetricData, ListMetrics, GetMetricStatistics, alarms                                            |
| `cloudformation-naming` | Auto physical name generation, explicit name precedence, cross-reference                           |
| `cognito`               | User pools, clients, AdminCreateUser, InitiateAuth, GetUser                                        |
| `cognito-oauth`         | Resource server CRUD, confidential clients, `/oauth2/token`, OIDC discovery, JWKS/JWT verification |
| `apigatewayv2`          | HTTP & WebSocket API lifecycle, routes, integrations, authorizers, stages, deployments, route responses, models, tagging |

## Requirements

- Node.js 20+
- npm

## Running

```bash
npm install

# All groups
npm test

# Via just (from compatibility-tests/)
just test-typescript
```

## Configuration

| Variable         | Default                 | Description             |
| ---------------- | ----------------------- | ----------------------- |
| `CLOUDSTACK_ENDPOINT` | `http://localhost:4566` | CloudStack emulator endpoint |

AWS credentials are always `test` / `test` / `us-east-1`.

## Docker

```bash
docker build -t CloudStack-sdk-node .
docker run --rm --network host CloudStack-sdk-node

# Custom endpoint (macOS/Windows)
docker run --rm -e CLOUDSTACK_ENDPOINT=http://host.docker.internal:4566 CloudStack-sdk-node
```

## TLS Tests

The `tls.test.ts` file tests HTTPS connectivity with AWS SDK v3 clients.
The `apigatewayv2-websocket-tls.test.ts` file tests WebSocket over TLS (WSS),
covering connect/disconnect, route selection, broadcast, authorization,
`@connections` API, binary frames, and payload limits.

These tests always run as part of the suite — they derive the HTTPS/WSS endpoint
from the HTTP endpoint by swapping the scheme (`http://CloudStack:4566` → `https://CloudStack:4566`,
`ws://` → `wss://`).

Since CloudStack serves HTTP and HTTPS simultaneously when TLS is enabled, existing HTTP
tests are unaffected.

The Dockerfile sets `NODE_TLS_REJECT_UNAUTHORIZED=0` so self-signed certs are accepted.
The CI workflow starts CloudStack with `CLOUDSTACK_TLS_ENABLED=true`.

To run locally:

```bash
# 1. Start CloudStack with TLS enabled
CLOUDSTACK_TLS_ENABLED=true ./mvnw quarkus:dev

# 2. Run all tests (including TLS)
NODE_TLS_REJECT_UNAUTHORIZED=0 npm test
```
