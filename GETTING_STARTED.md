# Getting Started with CloudStack

CloudStack is a lightweight, Go-native AWS emulator designed for local development and CI/CD. It supports over **65+ AWS services** out-of-the-box.

## 1. Installation
Ensure you have all the necessary tools installed (Go, AWS CLI, Just, JQ).

See the [Detailed Installation Guide](docs/INSTALLATION.md) for step-by-step instructions.

## 2. Building and Running

```bash
# Build
go build -o cloudstack ./cmd/cloudstack

# Run
./cloudstack
```
By default, it listens on port **4566**.

## 3. Configuring your Environment
You must redirect your AWS tools to the local emulator:

```bash
export AWS_ENDPOINT_URL=http://localhost:4566
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
```

*Tip: Use the `./bin/awslocal` wrapper to automate this configuration.*

## 4. Verification
Try creating an S3 bucket:
```bash
./bin/awslocal s3 mb s3://test-bucket
```

## 5. Testing
We provide a comprehensive integration test suite for all 65+ services.

See the [Testing Guide](docs/TESTING_GUIDE.md) to learn how to run them.
