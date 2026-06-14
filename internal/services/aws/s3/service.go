package s3

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/aws/s3/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type S3Service struct {
	bucketStore storage.Backend[string, *model.Bucket]
	objectStore storage.Backend[string, *model.S3Object]
	dataRoot    string
	inMemory    bool
	memoryData  map[string][]byte
	mu          sync.RWMutex
}

func NewS3Service(factory *storage.Factory, dataRoot string, inMemory bool) (*S3Service, error) {
	bucketStore, _ := storage.CreateAccountAware[*model.Bucket](factory, "s3", "s3-buckets.json", "wal")
	objectStore, _ := storage.CreateAccountAware[*model.S3Object](factory, "s3", "s3-objects.json", "wal")

	s := &S3Service{
		bucketStore: bucketStore,
		objectStore: objectStore,
		dataRoot:    dataRoot,
		inMemory:    inMemory,
		memoryData:  make(map[string][]byte),
	}

	if !inMemory {
		if err := os.MkdirAll(dataRoot, 0755); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *S3Service) CreateBucket(ctx context.Context, name, region string) (*model.Bucket, error) {
	if _, ok, _ := s.bucketStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("BucketAlreadyExists: The requested bucket name is not available")
	}

	bucket := &model.Bucket{
		Name:         name,
		CreationDate: time.Now(),
		Region:       region,
	}

	if err := s.bucketStore.Put(ctx, name, bucket); err != nil {
		return nil, err
	}

	if !s.inMemory {
		os.MkdirAll(filepath.Join(s.dataRoot, name), 0755)
	}

	log.Printf("Created S3 bucket: %s", name)
	return bucket, nil
}

func (s *S3Service) ListBuckets(ctx context.Context) ([]*model.Bucket, error) {
	return s.bucketStore.Scan(ctx, func(k string) bool { return true })
}

func (s *S3Service) PutObject(ctx context.Context, bucketName, key string, data []byte, contentType string, metadata map[string]string) (*model.S3Object, error) {
	if _, ok, _ := s.bucketStore.Get(ctx, bucketName); !ok {
		return nil, fmt.Errorf("NoSuchBucket: The specified bucket does not exist")
	}

	etag := fmt.Sprintf("\"%x\"", md5.Sum(data))
	obj := &model.S3Object{
		BucketName:   bucketName,
		Key:          key,
		ContentType:  contentType,
		Size:         int64(len(data)),
		LastModified: time.Now(),
		ETag:         etag,
		StorageClass: "STANDARD",
		IsLatest:     true,
		Metadata:     metadata,
	}

	// Store metadata
	if err := s.objectStore.Put(ctx, s.objectKey(bucketName, key), obj); err != nil {
		return nil, err
	}

	// Store data
	if s.inMemory {
		s.mu.Lock()
		s.memoryData[s.objectKey(bucketName, key)] = data
		s.mu.Unlock()
	} else {
		path := filepath.Join(s.dataRoot, bucketName, key)
		os.MkdirAll(filepath.Dir(path), 0755)
		if err := os.WriteFile(path, data, 0644); err != nil {
			return nil, err
		}
	}

	return obj, nil
}

func (s *S3Service) GetObject(ctx context.Context, bucketName, key string) (*model.S3Object, []byte, error) {
	obj, ok, err := s.objectStore.Get(ctx, s.objectKey(bucketName, key))
	if err != nil {
		return nil, nil, err
	}
	if !ok || obj.DeleteMarker {
		return nil, nil, fmt.Errorf("NoSuchKey: The specified key does not exist")
	}

	var data []byte
	if s.inMemory {
		s.mu.RLock()
		data = s.memoryData[s.objectKey(bucketName, key)]
		s.mu.RUnlock()
	} else {
		path := filepath.Join(s.dataRoot, bucketName, key)
		data, err = os.ReadFile(path)
		if err != nil {
			return nil, nil, err
		}
	}

	return obj, data, nil
}

func (s *S3Service) objectKey(bucket, key string) string {
	return fmt.Sprintf("%s/%s", bucket, key)
}
