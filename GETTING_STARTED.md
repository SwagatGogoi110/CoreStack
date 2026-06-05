# Getting Started with CloudStack

CloudStack is a lightweight, Go-native AWS emulator designed for local development and CI/CD.

## Prerequisites
- **Go 1.24+** (for local builds)
- **Docker** (for containerized execution)

## 1. Building from Source
To build the CloudStack binary locally:

```bash
go mod tidy
go build -o cloudstack ./cmd/cloudstack
```

## 2. Running CloudStack

### Local Binary
```bash
./cloudstack
```
By default, it listens on port **4566**.

### Using Docker
```bash
docker build -t cloudstack:latest .
docker run -p 4566:4566 -v /var/run/docker.sock:/var/run/docker.sock cloudstack:latest
```

### Using Docker Compose
Create a `docker-compose.yml`:
```yaml
services:
  cloudstack:
    image: cloudstack:latest
    ports:
      - "4566:4566"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ./data:/app/data
```
Run with `docker compose up -d`.

## 3. Configuring your Environment
Point your AWS CLI or SDK to the local endpoint:

```bash
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
export AWS_ENDPOINT_URL=http://localhost:4566
```

## 4. Testing the Connection
Verify that CloudStack is running:

```bash
curl http://localhost:4566/_cloudstack/health
```

Try creating an S3 bucket:
```bash
aws s3 mb s3://test-bucket --endpoint-url http://localhost:4566
```

## 5. Directory Structure
- `cmd/cloudstack`: Application entry point.
- `internal/core`: Dispatcher, protocol handlers, and routing.
- `internal/storage`: WAL and Memory storage backends.
- `internal/services`: Implementations for 50+ AWS services.
- `docs/`: Architectural details and service-specific guides.

## 6. Supported Services
CloudStack currently emulates over 50 AWS services, including:
- **Compute:** Lambda, EC2, ECS, EKS
- **Storage:** S3, DynamoDB
- **Messaging:** SQS, SNS, Kinesis, EventBridge
- **Security:** IAM, STS, Cognito, KMS, Secrets Manager
- **Analytics:** Athena, Glue, OpenSearch
