package s3

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/hectorvent/cloudstack/internal/storage"
)

func TestS3Service_CreateBucket(t *testing.T) {
	store := storage.NewFactory("./data-test", "000000000000", 1*time.Second)
	
	// Create a temporary directory for S3 data
	tempDir, err := os.MkdirTemp("", "s3-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	service, err := NewS3Service(store, tempDir, true)
	if err != nil {
		t.Fatalf("Failed to create S3 service: %v", err)
	}

	bucketName := "test-bucket"
	bucket, err := service.CreateBucket(context.Background(), bucketName, "us-east-1")
	if err != nil {
		t.Fatalf("Failed to create bucket: %v", err)
	}

	if bucket == nil {
		t.Fatalf("Bucket %s not returned after creation", bucketName)
	}
	if bucket.Name != bucketName {
		t.Errorf("Expected bucket name %s, got %s", bucketName, bucket.Name)
	}
}

func TestS3Service_PutAndGetObject(t *testing.T) {
	store := storage.NewFactory("./data-test", "000000000000", 1*time.Second)
	tempDir, _ := os.MkdirTemp("", "s3-test")
	defer os.RemoveAll(tempDir)

	service, _ := NewS3Service(store, tempDir, true)
	bucketName := "test-bucket"
	_, _ = service.CreateBucket(context.Background(), bucketName, "us-east-1")

	key := "test-object.txt"
	content := []byte("hello world from cloudstack s3")
	
	// Put Object
	objPut, err := service.PutObject(context.Background(), bucketName, key, content, "text/plain", nil)
	if err != nil {
		t.Fatalf("Failed to put object: %v", err)
	}

	if objPut.Size != int64(len(content)) {
		t.Errorf("Expected put object size %d, got %d", len(content), objPut.Size)
	}

	// Get Object
	obj, data, err := service.GetObject(context.Background(), bucketName, key)
	if err != nil {
		t.Fatalf("Failed to get object: %v", err)
	}

	if obj.Size != int64(len(content)) {
		t.Errorf("Expected object size %d, got %d", len(content), obj.Size)
	}

	if string(data) != string(content) {
		t.Errorf("Expected content %s, got %s", string(content), string(data))
	}
}
