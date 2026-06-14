package iam

import (
	"context"
	"testing"
	"time"

	"github.com/hectorvent/cloudstack/internal/storage"
)

func TestIamService_CreateUser(t *testing.T) {
	// Setup in-memory storage for testing
	store := storage.NewFactory("./data-test", "000000000000", 1*time.Second)
	service, err := NewIamService(store)
	if err != nil {
		t.Fatalf("Failed to create IAM service: %v", err)
	}

	userName := "test-user"
	user, err := service.CreateUser(context.Background(), userName, "/")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.UserName != userName {
		t.Errorf("Expected user name %s, got %s", userName, user.UserName)
	}
	if user.Arn == "" {
		t.Error("Expected ARN to be generated, got empty string")
	}
	if user.UserID == "" {
		t.Error("Expected UserId to be generated, got empty string")
	}
}

func TestIamService_CreateAccessKey(t *testing.T) {
	store := storage.NewFactory("./data-test", "000000000000", 1*time.Second)
	service, err := NewIamService(store)
	if err != nil {
		t.Fatalf("Failed to create IAM service: %v", err)
	}

	userName := "test-user"
	_, _ = service.CreateUser(context.Background(), userName, "/")

	key, err := service.CreateAccessKey(context.Background(), userName)
	if err != nil {
		t.Fatalf("Failed to create access key: %v", err)
	}

	if key.UserName != userName {
		t.Errorf("Expected access key user name %s, got %s", userName, key.UserName)
	}
	if len(key.AccessKeyID) < 16 {
		t.Errorf("Expected valid AccessKeyId, got %s", key.AccessKeyID)
	}
	if len(key.SecretAccessKey) < 32 {
		t.Errorf("Expected valid SecretAccessKey, got %s", key.SecretAccessKey)
	}
}
