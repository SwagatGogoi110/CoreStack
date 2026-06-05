# Running with Docker

CloudStack is distributed as a Docker image. All configuration is done through environment variables — no config files or volume-mounted YAML is required.

## Quick Start

```bash
docker run --rm -p 4566:4566 CloudStack/CloudStack:latest
```

That's it. The default configuration works out of the box for most services: SQS, SNS, S3, DynamoDB, SSM, Lambda, API Gateway, Cognito, KMS, Kinesis, Secrets Manager, CloudFormation, Step Functions, IAM, STS, EventBridge, Scheduler, and CloudWatch.

## Docker Compose

### Minimal (stateless)

```yaml title="docker-compose.yml"
services:
  CloudStack:
    image: CloudStack/CloudStack:latest
    ports:
      - "4566:4566"
    environment:
      CLOUDSTACK_HOSTNAME: CloudStack
```

### With persistence

Add two env vars and a volume — no config file needed:

```yaml title="docker-compose.yml"
services:
  CloudStack:
    image: CloudStack/CloudStack:latest
    ports:
      - "4566:4566"
    volumes:
      - CloudStack-data:/app/data
    environment:
      CLOUDSTACK_HOSTNAME: CloudStack
      CLOUDSTACK_STORAGE_MODE: persistent
      CLOUDSTACK_STORAGE_PERSISTENT_PATH: /app/data

volumes:
  CloudStack-data:
```

### With ElastiCache and RDS

ElastiCache and RDS proxy TCP connections to real Docker containers. Those containers' ports must be reachable from your host, so additional port ranges are exposed. The Docker socket is required for CloudStack to manage those containers:

```yaml title="docker-compose.yml"
services:
  CloudStack:
    image: CloudStack/CloudStack:latest
    ports:
      - "4566:4566"
      - "6379-6399:6379-6399"  # ElastiCache proxy ports
      - "7001-7099:7001-7099"  # RDS proxy ports
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - CloudStack-data:/app/data
    environment:
      CLOUDSTACK_HOSTNAME: CloudStack
      CLOUDSTACK_SERVICES_DOCKER_NETWORK: myproject_default
      CLOUDSTACK_STORAGE_MODE: persistent
      CLOUDSTACK_STORAGE_PERSISTENT_PATH: /app/data

volumes:
  CloudStack-data:
```

!!! warning "Docker socket"
    Lambda, ElastiCache, RDS, OpenSearch, and MSK require access to the Docker socket (`/var/run/docker.sock`) to spawn and manage containers. If you don't use these services, you can omit that volume.

!!! note "ECR port"
    ECR is backed by a `registry:2` sidecar container (`CloudStack-ecr-registry`) that CloudStack starts and manages. That container binds its own host port (default `5100`) directly — do not add `5100-5199` to the CloudStack service's `ports` list. See [Ports Reference → ECR](./ports.md#ports-51005199--ecr-registry).

## Multi-container Networking

By default CloudStack embeds `localhost` in response URLs — for example, SQS queue URLs look like `http://localhost:4566/000000000000/my-queue`. This works when your application runs on the same machine, but breaks inside Docker Compose because other containers cannot reach `localhost` of the CloudStack container.

Set `CLOUDSTACK_HOSTNAME` to the Compose service name so CloudStack uses that name in every URL it generates:

```yaml title="docker-compose.yml"
services:
  CloudStack:
    image: CloudStack/CloudStack:latest
    ports:
      - "4566:4566"
    environment:
      CLOUDSTACK_HOSTNAME: CloudStack   # (1)

  app:
    build: .
    environment:
      AWS_ENDPOINT_URL: http://CloudStack:4566
    depends_on:
      - CloudStack
```

1. Must match the Compose service name so other containers can resolve it by DNS.

With this setting CloudStack returns URLs like `http://CloudStack:4566/000000000000/my-queue` that other containers can reach.

CloudStack automatically attaches Lambda containers it spawns to the same Compose network when no explicit Docker network is configured. Setting `CLOUDSTACK_HOSTNAME` ensures those containers receive a reachable endpoint, and that response fields such as SQS `QueueUrl` use the Docker service name instead of `localhost`.

Fields affected:

- SQS — `QueueUrl`
- SNS — topic ARN callback URLs and subscription endpoints
- Any pre-signed URL or callback generated from `CLOUDSTACK_BASE_URL`

!!! tip "CI pipelines"
    In GitHub Actions or GitLab CI where both your app and CloudStack run as `services`, set `CLOUDSTACK_HOSTNAME` to the service name (e.g. `CloudStack`) and point your SDK at `http://CloudStack:4566`.

## Initialization Hooks

Mount shell scripts into hook directories to run setup or teardown logic at each lifecycle phase. No configuration variable is needed — CloudStack detects scripts by directory:

```yaml title="docker-compose.yml"
services:
  CloudStack:
    image: CloudStack/CloudStack:latest-compat
    ports:
      - "4566:4566"
    volumes:
      - ./init/boot.d:/etc/CloudStack/init/boot.d:ro    # before storage loads — no AWS APIs yet
      - ./init/start.d:/etc/CloudStack/init/start.d:ro  # after HTTP server is ready
      - ./init/ready.d:/etc/CloudStack/init/ready.d:ro  # after all start hooks complete
      - ./init/stop.d:/etc/CloudStack/init/stop.d:ro    # during shutdown, while HTTP is still up
```

Use the `latest-compat` image when your scripts call `aws` or `boto3` — it includes the AWS CLI and boto3 pre-configured for the local endpoint, so no `--endpoint-url` flag is needed.

If you have existing LocalStack init scripts, mount them under the LocalStack-compat paths and they run unchanged:

```yaml
volumes:
  - ./localstack-init/ready.d:/etc/localstack/init/ready.d:ro
```

See [Initialization Hooks](./initialization-hooks.md) for execution order, script types, and exit-code behavior.

## CI Pipeline Example

```yaml title=".github/workflows/test.yml"
services:
  CloudStack:
    image: CloudStack/CloudStack:latest
    ports:
      - "4566:4566"

steps:
  - name: Run tests
    env:
      AWS_ENDPOINT_URL: http://localhost:4566
      AWS_DEFAULT_REGION: us-east-1
      AWS_ACCESS_KEY_ID: test
      AWS_SECRET_ACCESS_KEY: test
    run: mvn test
```

## Common Environment Variables

The most frequently set variables when running CloudStack as a Docker image:

| Variable | Default | Purpose |
|---|---|---|
| `CLOUDSTACK_HOSTNAME` | _(none)_ | Hostname embedded in response URLs. Set to the Compose service name in multi-container setups |
| `CLOUDSTACK_DEFAULT_REGION` | `us-east-1` | AWS region reported in ARNs and responses |
| `CLOUDSTACK_DEFAULT_ACCOUNT_ID` | `000000000000` | AWS account ID used in ARNs |
| `CLOUDSTACK_STORAGE_MODE` | `memory` | `memory`, `persistent`, `hybrid`, or `wal` |
| `CLOUDSTACK_STORAGE_PERSISTENT_PATH` | `./data` | Directory for persistent storage |
| `CLOUDSTACK_SERVICES_DOCKER_NETWORK` | _(none)_ | Docker network for spawned containers (Lambda, ElastiCache, RDS, OpenSearch, MSK) |
| `CLOUDSTACK_AUTH_VALIDATE_SIGNATURES` | `false` | Verify AWS request signatures |
| `CLOUDSTACK_SERVICES_LAMBDA_EPHEMERAL` | `false` | Remove Lambda containers after each invocation |

For the complete list of every `CLOUDSTACK_*` variable, see [Environment Variables Reference](./environment-variables.md).

## Docker Configuration

For Docker daemon socket, private registry authentication, log rotation, and network settings see [Docker Configuration](./docker.md).
