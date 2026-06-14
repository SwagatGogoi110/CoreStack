# Installation & Requirements

This document outlines the necessary tools and steps to get CloudStack running and tested on your local machine.

## Prerequisites

### 1. Go (Compiler)
CloudStack is written in Go. You need Go 1.24+ installed.
*   **macOS:** `brew install go`
*   **Linux/Windows:** [Download from go.dev](https://go.dev/dl/)

### 2. AWS CLI (Client)
The standard AWS Command Line Interface is used to interact with CloudStack.
*   **macOS:** `brew install awscli`
*   **Linux/Windows:** [Install AWS CLI v2](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)

### 3. Just (Task Runner)
We use `just` to automate testing and setup.
*   **macOS:** `brew install just`
*   **Linux:** `sudo apt install just` or `cargo install just`

### 4. JQ (JSON Processor)
Essential for the integration test suite.
*   **macOS:** `brew install jq`
*   **Linux:** `sudo apt install jq`

---

## Setting up the Environment

### 1. Clone the Repository
```bash
git clone https://github.com/hectorvent/cloudstack.git
cd cloudstack
```

### 2. Build the Application
```bash
go build -o cloudstack ./cmd/cloudstack
```

### 3. Install Test Dependencies
We use BATS for integration testing. You can install it using the provided `just` recipe:
```bash
cd compatibility-tests
just setup-awscli
```

---

## Configuration

CloudStack is configured via environment variables. The most important one is `AWS_ENDPOINT_URL`, which redirects CLI commands to your local machine.

### Automatic Configuration (The `awslocal` Wrapper)
We provide a convenience script in `bin/awslocal` that automatically handles the endpoint redirection for you.

```bash
# Instead of:
# aws s3 ls --endpoint-url http://localhost:4566

# You can just use:
./bin/awslocal s3 ls
```

### Manual Configuration
If you prefer using the standard `aws` command, export these variables in your terminal:

```bash
export AWS_ENDPOINT_URL=http://localhost:4566
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
```

Now you are ready to start [Testing](TESTING_GUIDE.md)!
