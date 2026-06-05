# Environment Variables Reference

CloudStack is configured exclusively through environment variables. Every option below maps directly to a `CLOUDSTACK_*` variable — no YAML file is needed when running the published Docker image.

---

## Global

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_BASE_URL` | `http://localhost:4566` | Base URL embedded in response fields (SQS `QueueUrl`, pre-signed URLs, etc.) |
| `CLOUDSTACK_HOSTNAME` | _(none)_ | Overrides only the hostname part of `CLOUDSTACK_BASE_URL`. Set to the Compose service name (e.g. `CloudStack`) so other containers can reach CloudStack by DNS |
| `CLOUDSTACK_DEFAULT_REGION` | `us-east-1` | AWS region used in ARNs and API responses |
| `CLOUDSTACK_DEFAULT_ACCOUNT_ID` | `000000000000` | Fallback account ID used in ARNs when the request's access key is not exactly 12 digits. When the access key IS 12 digits, it is used directly as the account ID — see [Multi-Account Isolation](./multi-account.md) |
| `CLOUDSTACK_DEFAULT_AVAILABILITY_ZONE` | `us-east-1a` | Availability zone reported in EC2 and other responses |
| `CLOUDSTACK_MAX_REQUEST_SIZE` | `512` | Maximum HTTP request body size in megabytes |
| `CLOUDSTACK_ECR_BASE_URI` | `public.ecr.aws` | Base URI for public ECR image references |

---

## Authentication

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_AUTH_VALIDATE_SIGNATURES` | `false` | When `true`, verifies AWS Signature V4 on every request. Leave `false` for local development |
| `CLOUDSTACK_AUTH_PRESIGN_SECRET` | `local-emulator-secret` | Secret used to sign and verify pre-signed URLs |

## Browser CORS

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SECURITY_EXTRA_CORS_ALLOWED_ORIGINS` | _(none)_ | Comma-separated browser origins allowed to call CloudStack directly. Alias: `EXTRA_CORS_ALLOWED_ORIGINS` |
| `CLOUDSTACK_SECURITY_EXTRA_CORS_ALLOWED_HEADERS` | _(none)_ | Additional header names to include in `Access-Control-Allow-Headers`. Alias: `EXTRA_CORS_ALLOWED_HEADERS` |
| `CLOUDSTACK_SECURITY_EXTRA_CORS_EXPOSE_HEADERS` | _(none)_ | Additional header names to include in `Access-Control-Expose-Headers`. Alias: `EXTRA_CORS_EXPOSE_HEADERS` |
| `CLOUDSTACK_SECURITY_DISABLE_CORS_HEADERS` | `false` | Disable CloudStack's global CORS response headers. Alias: `DISABLE_CORS_HEADERS` |

---

## TLS / HTTPS

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_TLS_ENABLED` | `false` | Enable TLS/HTTPS on all endpoints (HTTP remains available simultaneously) |
| `CLOUDSTACK_TLS_CERT_PATH` | _(none)_ | Path to a PEM certificate file. When set, disables auto-generation |
| `CLOUDSTACK_TLS_KEY_PATH` | _(none)_ | Path to a PEM private key file. Required when `CLOUDSTACK_TLS_CERT_PATH` is set |
| `CLOUDSTACK_TLS_SELF_SIGNED` | `true` | Auto-generate and persist a self-signed certificate when no cert/key paths are provided |

See [TLS / HTTPS](./tls.md) for SDK configuration examples and WebSocket (`wss://`) support.

---

## Storage

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_STORAGE_MODE` | `memory` | Global storage backend: `memory`, `persistent`, `hybrid`, or `wal` |
| `CLOUDSTACK_STORAGE_PERSISTENT_PATH` | `./data` | Container-side directory for persistent and hybrid storage |
| `CLOUDSTACK_STORAGE_HOST_PERSISTENT_PATH` | `./data` | Host-side path for Docker volume bind-mounts (RDS, OpenSearch, MSK, ECR data). When unset, CloudStack uses named Docker volumes |
| `CLOUDSTACK_STORAGE_PRUNE_VOLUMES_ON_DELETE` | `false` | Remove named Docker volumes immediately when the resource is deleted |
| `CLOUDSTACK_STORAGE_WAL_COMPACTION_INTERVAL_MS` | `30000` | How often (ms) the WAL compaction runs. Only applies when `CLOUDSTACK_STORAGE_MODE=wal` |

### Per-service storage overrides

Each service can override the global storage mode and flush interval. Replace `<SERVICE>` with the uppercase service name:

```
CLOUDSTACK_STORAGE_SERVICES_<SERVICE>_MODE=hybrid
CLOUDSTACK_STORAGE_SERVICES_<SERVICE>_FLUSH_INTERVAL_MS=5000
```

Available service names: `SSM`, `SQS`, `S3`, `DYNAMODB`, `SNS`, `LAMBDA`, `CLOUDWATCHLOGS`, `CLOUDWATCHMETRICS`, `SECRETSMANAGER`, `ACM`, `OPENSEARCH`, `RDS`, `ELASTICACHE`, `APPCONFIG`, `APPCONFIGDATA`, `BACKUP`.

See [Storage Modes](./storage.md) for a full explanation of each mode.

---

## Docker Daemon

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_DOCKER_DOCKER_HOST` | `unix:///var/run/docker.sock` | Docker daemon socket path or TCP address |
| `CLOUDSTACK_DOCKER_DOCKER_CONFIG_PATH` | _(none)_ | Path to a directory containing Docker's `config.json` for registry auth |
| `CLOUDSTACK_DOCKER_LOG_MAX_SIZE` | `10m` | Log rotation max size for spawned containers (e.g. `10m`, `1g`) |
| `CLOUDSTACK_DOCKER_LOG_MAX_FILE` | `3` | Number of rotated log files to keep for spawned containers |

### Registry credentials

Provide credentials for private registries (e.g. for Lambda base images). Use incrementing indexes (`0`, `1`, `2`, …) for multiple registries:

| Variable | Description |
|---|---|
| `CLOUDSTACK_DOCKER_REGISTRY_CREDENTIALS_0__SERVER` | Registry hostname (e.g. `ghcr.io`) |
| `CLOUDSTACK_DOCKER_REGISTRY_CREDENTIALS_0__USERNAME` | Registry username |
| `CLOUDSTACK_DOCKER_REGISTRY_CREDENTIALS_0__PASSWORD` | Registry password or token |

---

## DNS

CloudStack's embedded DNS server always resolves the following wildcard suffixes to CloudStack's container IP — no configuration required:

| Built-in suffix | Covers |
|---|---|
| `localhost.cloudstack.io` | `localhost.cloudstack.io` and `*.localhost.cloudstack.io` (e.g. `my-bucket.s3.localhost.cloudstack.io`) |
| `localhost.localstack.cloud` | `localhost.localstack.cloud` and `*.localhost.localstack.cloud` (e.g. `my-bucket.s3.localhost.localstack.cloud`) |

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_DNS_EXTRA_SUFFIXES` | _(none)_ | Comma-separated list of additional hostname suffixes to resolve to CloudStack's container IP. Use this for custom domains beyond the built-in ones above (e.g. a private internal suffix). |

---

## Initialization Hooks

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_INIT_HOOKS_SHELL_EXECUTABLE` | `/bin/sh` | Shell used to execute hook scripts |
| `CLOUDSTACK_INIT_HOOKS_TIMEOUT_SECONDS` | `30` | Maximum time a single hook script may run |
| `CLOUDSTACK_INIT_HOOKS_SHUTDOWN_GRACE_PERIOD_SECONDS` | `2` | Time allowed for stop hooks to complete during shutdown |

See [Initialization Hooks](./initialization-hooks.md) for lifecycle phases and script conventions.

---

## Services — Shared

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_DOCKER_NETWORK` | _(none)_ | Docker network name used by all container-backed services (Lambda, RDS, ElastiCache, ECS, OpenSearch, EKS, MSK). Per-service overrides take precedence |

---

## Services — Core

### SSM (Parameter Store)

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_SSM_ENABLED` | `true` | Enable the SSM service |
| `CLOUDSTACK_SERVICES_SSM_MAX_PARAMETER_HISTORY` | `5` | Maximum number of historical versions kept per parameter |

### SQS

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_SQS_ENABLED` | `true` | Enable the SQS service |
| `CLOUDSTACK_SERVICES_SQS_DEFAULT_VISIBILITY_TIMEOUT` | `30` | Default message visibility timeout in seconds |
| `CLOUDSTACK_SERVICES_SQS_MAX_MESSAGE_SIZE` | `262144` | Maximum message body size in bytes (256 KB) |
| `CLOUDSTACK_SERVICES_SQS_CLEAR_FIFO_DEDUPLICATION_CACHE_ON_PURGE` | `false` | Reset the deduplication cache when a FIFO queue is purged |

### SNS

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_SNS_ENABLED` | `true` | Enable the SNS service |

### S3

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_S3_ENABLED` | `true` | Enable the S3 service |
| `CLOUDSTACK_SERVICES_S3_DEFAULT_PRESIGN_EXPIRY_SECONDS` | `3600` | Default pre-signed URL expiry when none is specified |

### DynamoDB

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_DYNAMODB_ENABLED` | `true` | Enable the DynamoDB service |

### Lambda

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_LAMBDA_ENABLED` | `true` | Enable the Lambda service |
| `CLOUDSTACK_SERVICES_LAMBDA_EPHEMERAL` | `false` | Remove Lambda containers immediately after each invocation |
| `CLOUDSTACK_SERVICES_LAMBDA_DEFAULT_MEMORY_MB` | `128` | Default memory allocation for functions that don't specify one |
| `CLOUDSTACK_SERVICES_LAMBDA_DEFAULT_TIMEOUT_SECONDS` | `3` | Default invocation timeout in seconds |
| `CLOUDSTACK_SERVICES_LAMBDA_RUNTIME_API_BASE_PORT` | `9200` | First port in the Lambda Runtime API port range |
| `CLOUDSTACK_SERVICES_LAMBDA_RUNTIME_API_MAX_PORT` | `9299` | Last port in the Lambda Runtime API port range |
| `CLOUDSTACK_SERVICES_LAMBDA_CODE_PATH` | `./data/lambda-code` | Container path where Lambda deployment ZIPs are stored |
| `CLOUDSTACK_SERVICES_LAMBDA_POLL_INTERVAL_MS` | `1000` | How often (ms) the SQS and Kinesis event source pollers check for new messages |
| `CLOUDSTACK_SERVICES_LAMBDA_CONTAINER_IDLE_TIMEOUT_SECONDS` | `300` | Seconds of inactivity before an idle Lambda container is removed |
| `CLOUDSTACK_SERVICES_LAMBDA_REGION_CONCURRENCY_LIMIT` | `1000` | Maximum concurrent Lambda invocations across all functions in a region |
| `CLOUDSTACK_SERVICES_LAMBDA_UNRESERVED_CONCURRENCY_MIN` | `100` | Minimum unreserved concurrency pool |
| `CLOUDSTACK_SERVICES_LAMBDA_HOT_RELOAD_ENABLED` | `false` | Watch Lambda code directories for changes and reload without redeployment |
| `CLOUDSTACK_SERVICES_LAMBDA_HOT_RELOAD_ALLOWED_PATHS` | _(none)_ | Comma-separated host paths that hot-reload is allowed to watch |
| `CLOUDSTACK_SERVICES_LAMBDA_DOCKER_NETWORK` | _(none)_ | Docker network for Lambda containers (overrides `CLOUDSTACK_SERVICES_DOCKER_NETWORK`) |
| `CLOUDSTACK_SERVICES_LAMBDA_AWS_CONFIG_PATH` | _(none)_ | Host path bind-mounted read-only at `/opt/aws-config` inside Lambda containers for real credential discovery |

### API Gateway

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_APIGATEWAY_ENABLED` | `true` | Enable the API Gateway v1 (REST) service |
| `CLOUDSTACK_SERVICES_APIGATEWAYV2_ENABLED` | `true` | Enable the API Gateway v2 (HTTP + WebSocket) service |

### IAM

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_IAM_ENABLED` | `true` | Enable the IAM service |
| `CLOUDSTACK_SERVICES_IAM_ENFORCEMENT_ENABLED` | `false` | When `true`, enforce IAM policies on API calls. Leave `false` for most local development scenarios |

### KMS

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_KMS_ENABLED` | `true` | Enable the KMS service |

### Kinesis

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_KINESIS_ENABLED` | `true` | Enable the Kinesis Data Streams service |

### Firehose

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_FIREHOSE_ENABLED` | `true` | Enable the Kinesis Data Firehose service |

### EventBridge

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_EVENTBRIDGE_ENABLED` | `true` | Enable the EventBridge service |

### Scheduler

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_SCHEDULER_ENABLED` | `true` | Enable the EventBridge Scheduler service |
| `CLOUDSTACK_SERVICES_SCHEDULER_INVOCATION_ENABLED` | `true` | When `false`, schedules are stored but never invoked |
| `CLOUDSTACK_SERVICES_SCHEDULER_TICK_INTERVAL_SECONDS` | `10` | How often (seconds) the scheduler checks for due schedules |

### CloudWatch Logs

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_CLOUDWATCHLOGS_ENABLED` | `true` | Enable the CloudWatch Logs service |
| `CLOUDSTACK_SERVICES_CLOUDWATCHLOGS_MAX_EVENTS_PER_QUERY` | `10000` | Maximum log events returned by a single `FilterLogEvents` call |

### CloudWatch Metrics

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_CLOUDWATCHMETRICS_ENABLED` | `true` | Enable the CloudWatch Metrics service |

### Secrets Manager

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_SECRETSMANAGER_ENABLED` | `true` | Enable the Secrets Manager service |
| `CLOUDSTACK_SERVICES_SECRETSMANAGER_DEFAULT_RECOVERY_WINDOW_DAYS` | `30` | Default recovery window for deleted secrets |

### Cognito

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_COGNITO_ENABLED` | `true` | Enable the Cognito User Pools service |

### Step Functions

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_STEPFUNCTIONS_ENABLED` | `true` | Enable the Step Functions service |

### CloudFormation

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_CLOUDFORMATION_ENABLED` | `true` | Enable the CloudFormation service |

### ACM (Certificate Manager)

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_ACM_ENABLED` | `true` | Enable the ACM service |
| `CLOUDSTACK_SERVICES_ACM_VALIDATION_WAIT_SECONDS` | `0` | Simulated delay before a requested certificate transitions to `ISSUED` |

### SES (Simple Email Service)

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_SES_ENABLED` | `true` | Enable the SES service |
| `CLOUDSTACK_SERVICES_SES_SMTP_HOST` | _(none)_ | SMTP relay host for outbound email. When unset, emails are captured in memory only |
| `CLOUDSTACK_SERVICES_SES_SMTP_PORT` | `25` | SMTP relay port |
| `CLOUDSTACK_SERVICES_SES_SMTP_USER` | _(none)_ | SMTP username |
| `CLOUDSTACK_SERVICES_SES_SMTP_PASS` | _(none)_ | SMTP password |
| `CLOUDSTACK_SERVICES_SES_SMTP_STARTTLS` | `DISABLED` | STARTTLS mode: `DISABLED`, `OPTIONAL`, or `REQUIRED` |

### Pipes

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_PIPES_ENABLED` | `true` | Enable the EventBridge Pipes service |

---

## Services — Container-Backed

These services spawn Docker containers. They require access to the Docker socket (`/var/run/docker.sock`).

### ElastiCache

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_ELASTICACHE_ENABLED` | `true` | Enable the ElastiCache service |
| `CLOUDSTACK_SERVICES_ELASTICACHE_PROXY_BASE_PORT` | `6379` | First port in the ElastiCache proxy range |
| `CLOUDSTACK_SERVICES_ELASTICACHE_PROXY_MAX_PORT` | `6399` | Last port in the ElastiCache proxy range |
| `CLOUDSTACK_SERVICES_ELASTICACHE_DEFAULT_IMAGE` | `valkey/valkey:8` | Default Docker image for cache clusters |
| `CLOUDSTACK_SERVICES_ELASTICACHE_DOCKER_NETWORK` | _(none)_ | Docker network for ElastiCache containers (overrides `CLOUDSTACK_SERVICES_DOCKER_NETWORK`) |

### RDS

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_RDS_ENABLED` | `true` | Enable the RDS service |
| `CLOUDSTACK_SERVICES_RDS_PROXY_BASE_PORT` | `7001` | First port in the RDS proxy range |
| `CLOUDSTACK_SERVICES_RDS_PROXY_MAX_PORT` | `7099` | Last port in the RDS proxy range |
| `CLOUDSTACK_SERVICES_RDS_DEFAULT_POSTGRES_IMAGE` | `postgres:16-alpine` | Default PostgreSQL Docker image |
| `CLOUDSTACK_SERVICES_RDS_DEFAULT_MYSQL_IMAGE` | `mysql:8.0` | Default MySQL Docker image |
| `CLOUDSTACK_SERVICES_RDS_DEFAULT_MARIADB_IMAGE` | `mariadb:11` | Default MariaDB Docker image |
| `CLOUDSTACK_SERVICES_RDS_DOCKER_NETWORK` | _(none)_ | Docker network for RDS containers (overrides `CLOUDSTACK_SERVICES_DOCKER_NETWORK`) |

### OpenSearch

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_OPENSEARCH_ENABLED` | `true` | Enable the OpenSearch service |
| `CLOUDSTACK_SERVICES_OPENSEARCH_MOCK` | `false` | When `true`, domains are created instantly without a real container (API only) |
| `CLOUDSTACK_SERVICES_OPENSEARCH_DEFAULT_IMAGE` | `opensearchproject/opensearch:2` | Docker image for OpenSearch domains |
| `CLOUDSTACK_SERVICES_OPENSEARCH_PROXY_BASE_PORT` | `9400` | First port in the OpenSearch proxy range |
| `CLOUDSTACK_SERVICES_OPENSEARCH_PROXY_MAX_PORT` | `9499` | Last port in the OpenSearch proxy range |
| `CLOUDSTACK_SERVICES_OPENSEARCH_KEEP_RUNNING_ON_SHUTDOWN` | `false` | Keep OpenSearch containers running when CloudStack stops |

### MSK (Managed Streaming for Kafka)

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_MSK_ENABLED` | `true` | Enable the MSK service |
| `CLOUDSTACK_SERVICES_MSK_MOCK` | `false` | When `true`, clusters are created instantly without a real Redpanda container |
| `CLOUDSTACK_SERVICES_MSK_DEFAULT_IMAGE` | `redpandadata/redpanda:latest` | Docker image for Kafka/Redpanda brokers |

### ECR (Elastic Container Registry)

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_ECR_ENABLED` | `true` | Enable the ECR service |
| `CLOUDSTACK_SERVICES_ECR_REGISTRY_IMAGE` | `registry:2` | Docker image for the ECR registry sidecar |
| `CLOUDSTACK_SERVICES_ECR_REGISTRY_CONTAINER_NAME` | `CloudStack-ecr-registry` | Name of the ECR registry sidecar container |
| `CLOUDSTACK_SERVICES_ECR_REGISTRY_BASE_PORT` | `5100` | First port in the ECR registry range |
| `CLOUDSTACK_SERVICES_ECR_REGISTRY_MAX_PORT` | `5199` | Last port in the ECR registry range |
| `CLOUDSTACK_SERVICES_ECR_TLS_ENABLED` | `false` | Enable TLS for the ECR registry |
| `CLOUDSTACK_SERVICES_ECR_KEEP_RUNNING_ON_SHUTDOWN` | `true` | Keep the ECR registry container running when CloudStack stops |
| `CLOUDSTACK_SERVICES_ECR_URI_STYLE` | `hostname` | Repository URI style: `hostname` (`<account>.dkr.ecr.<region>.localhost`) or `path` |
| `CLOUDSTACK_SERVICES_ECR_DOCKER_NETWORK` | _(none)_ | Docker network for the ECR registry container |

### EKS (Elastic Kubernetes Service)

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_EKS_ENABLED` | `true` | Enable the EKS service |
| `CLOUDSTACK_SERVICES_EKS_MOCK` | `false` | When `true`, clusters are created instantly without a real container |
| `CLOUDSTACK_SERVICES_EKS_PROVIDER` | `k3s` | Kubernetes provider (`k3s`) |
| `CLOUDSTACK_SERVICES_EKS_DEFAULT_IMAGE` | `rancher/k3s:latest` | Docker image for EKS clusters |
| `CLOUDSTACK_SERVICES_EKS_API_SERVER_BASE_PORT` | `6500` | First port in the Kubernetes API server range |
| `CLOUDSTACK_SERVICES_EKS_API_SERVER_MAX_PORT` | `6599` | Last port in the Kubernetes API server range |
| `CLOUDSTACK_SERVICES_EKS_KEEP_RUNNING_ON_SHUTDOWN` | `false` | Keep EKS containers running when CloudStack stops |
| `CLOUDSTACK_SERVICES_EKS_DOCKER_NETWORK` | _(none)_ | Docker network for EKS containers |

### ECS (Elastic Container Service)

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_ECS_ENABLED` | `true` | Enable the ECS service |
| `CLOUDSTACK_SERVICES_ECS_MOCK` | `false` | When `true`, tasks are registered but not actually run |
| `CLOUDSTACK_SERVICES_ECS_DEFAULT_MEMORY_MB` | `512` | Default task memory when not specified in the task definition |
| `CLOUDSTACK_SERVICES_ECS_DEFAULT_CPU_UNITS` | `256` | Default task CPU units when not specified in the task definition |
| `CLOUDSTACK_SERVICES_ECS_DOCKER_NETWORK` | _(none)_ | Docker network for ECS task containers |

### EC2

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_EC2_ENABLED` | `true` | Enable the EC2 service |
| `CLOUDSTACK_SERVICES_EC2_MOCK` | `false` | When `true`, instances are registered in state but no containers are spawned |
| `CLOUDSTACK_SERVICES_EC2_IMDS_PORT` | `9169` | Port for the EC2 Instance Metadata Service (IMDS) endpoint |
| `CLOUDSTACK_SERVICES_EC2_SSH_PORT_RANGE_START` | `2200` | First port in the SSH port range for EC2 instances |
| `CLOUDSTACK_SERVICES_EC2_SSH_PORT_RANGE_END` | `2299` | Last port in the SSH port range |

### Athena

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_ATHENA_ENABLED` | `true` | Enable the Athena service |
| `CLOUDSTACK_SERVICES_ATHENA_MOCK` | `false` | When `true`, queries are accepted but not executed |
| `CLOUDSTACK_SERVICES_ATHENA_DEFAULT_IMAGE` | `CloudStack/CloudStack-duck:latest` | Docker image for the DuckDB query engine |
| `CLOUDSTACK_SERVICES_ATHENA_DUCK_URL` | _(none)_ | URL of an existing DuckDB service. When set, CloudStack skips managing the container |

---

## Services — Additional

| Variable | Default | Description |
|---|---|---|
| `CLOUDSTACK_SERVICES_GLUE_ENABLED` | `true` | Enable the Glue service |
| `CLOUDSTACK_SERVICES_BEDROCK_RUNTIME_ENABLED` | `true` | Enable the Bedrock Runtime service |
| `CLOUDSTACK_SERVICES_TEXTRACT_ENABLED` | `true` | Enable the Textract service |
| `CLOUDSTACK_SERVICES_TRANSFER_ENABLED` | `true` | Enable the Transfer Family service |
| `CLOUDSTACK_SERVICES_ROUTE53_ENABLED` | `true` | Enable the Route 53 service |
| `CLOUDSTACK_SERVICES_ELBV2_ENABLED` | `true` | Enable the ELBv2 (ALB/NLB) service |
| `CLOUDSTACK_SERVICES_ELBV2_MOCK` | `false` | When `true`, load balancers are registered but no containers are spawned |
| `CLOUDSTACK_SERVICES_AUTOSCALING_ENABLED` | `true` | Enable the Auto Scaling service |
| `CLOUDSTACK_SERVICES_CODEBUILD_ENABLED` | `true` | Enable the CodeBuild service |
| `CLOUDSTACK_SERVICES_CODEBUILD_DOCKER_NETWORK` | _(none)_ | Docker network for CodeBuild build containers |
| `CLOUDSTACK_SERVICES_CODEDEPLOY_ENABLED` | `true` | Enable the CodeDeploy service |
| `CLOUDSTACK_SERVICES_BACKUP_ENABLED` | `true` | Enable the AWS Backup service |
| `CLOUDSTACK_SERVICES_BACKUP_JOB_COMPLETION_DELAY_SECONDS` | `3` | Simulated delay before backup jobs transition to `COMPLETED` |
| `CLOUDSTACK_SERVICES_APPCONFIG_ENABLED` | `true` | Enable the AppConfig service |
| `CLOUDSTACK_SERVICES_APPCONFIGDATA_ENABLED` | `true` | Enable the AppConfig Data service |
