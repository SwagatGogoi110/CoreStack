# CloudStack Unit Testing Guide

This guide explains how to write, run, and structure unit tests for CloudStack services.

## Overview
CloudStack uses standard Go testing (`testing` package). We focus on testing the **Service Layer** (`service_test.go`), which contains the core business logic (e.g., entity creation, ID generation, storage interactions). 

We use an in-memory or temporary directory configuration for the `storage.Factory` so that unit tests run blazingly fast without leaving persistent files behind.

## Running Tests

We provide a convenience script: `test-service.sh` at the root of the repository.

### Run tests for a specific service:
```bash
./test-service.sh s3
./test-service.sh iam
./test-service.sh dynamodb
```

### Run all tests across all services:
```bash
./test-service.sh all
```
*(Or use the standard `go test -v ./internal/...`)*

## Writing Tests

Unit tests should be placed in `internal/services/<service_name>/service_test.go`.

### Basic Test Pattern
1. **Initialize Storage**: Create a test storage factory pointing to a dummy directory.
2. **Initialize Service**: Inject the storage factory into the service constructor.
3. **Execute Logic**: Call the service method.
4. **Assert**: Verify the returned object or error state.

### Example (IAM Service)
```go
package iam

import (
	"testing"
	"time"
	"github.com/hectorvent/cloudstack/internal/storage"
)

func TestIamService_CreateUser(t *testing.T) {
	// 1. Setup ephemeral storage
	store := storage.NewFactory("./data-test", "000000000000", 1*time.Second)
	
	// 2. Initialize Service
	service, _ := NewIamService(store)

	// 3. Execute
	user, err := service.CreateUser("test-user", "/")
	
	// 4. Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if user.UserName != "test-user" {
		t.Errorf("Expected test-user, got %v", user.UserName)
	}
}
```

## Adding Tests for New Services
When porting or updating a service, ensure that at minimum the basic CRUD operations (Create/Get) are covered by unit tests. If a service interacts with the file system (like S3 or Lambda), ensure you use `os.MkdirTemp` to create a temporary directory for the test duration, and `defer os.RemoveAll(tempDir)` to clean it up.
