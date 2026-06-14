package gcs

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/gcp/gcs/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type GcsService struct {
	bucketStore storage.Backend[string, *model.Bucket]
	objectStore storage.Backend[string, *model.ObjectMeta]
	dataRoot    string
	mu          sync.RWMutex
}

func NewGcsService(factory *storage.Factory) (*GcsService, error) {
	bucketStore, _ := storage.CreateAccountAware[*model.Bucket](factory, "gcs", "gcs-buckets.json", "wal")
	objectStore, _ := storage.CreateAccountAware[*model.ObjectMeta](factory, "gcs", "gcs-objects.json", "wal")

	dataRoot := filepath.Join(factory.Root(), "gcs-data")
	os.MkdirAll(dataRoot, 0755)

	return &GcsService{
		bucketStore: bucketStore,
		objectStore: objectStore,
		dataRoot:    dataRoot,
	}, nil
}

// Buckets

func (s *GcsService) CreateBucket(ctx context.Context, name string) (*model.Bucket, error) {
	if _, ok, _ := s.bucketStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("Conflict: Bucket %s already exists", name)
	}

	bucket := &model.Bucket{
		Kind:           "storage#bucket",
		Id:             name,
		Name:           name,
		SelfLink:       fmt.Sprintf("https://www.googleapis.com/storage/v1/b/%s", name),
		Location:       "US",
		StorageClass:   "STANDARD",
		TimeCreated:    time.Now(),
		Updated:        time.Now(),
		Metageneration: "1",
		Etag:           "CAE=",
	}

	if err := s.bucketStore.Put(ctx, name, bucket); err != nil {
		return nil, err
	}

	// Create physical directory
	os.MkdirAll(filepath.Join(s.dataRoot, name), 0755)

	return bucket, nil
}

func (s *GcsService) GetBucket(ctx context.Context, name string) (*model.Bucket, error) {
	b, ok, err := s.bucketStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NotFound: Bucket %s not found", name)
	}
	return b, nil
}

func (s *GcsService) ListBuckets(ctx context.Context) ([]*model.Bucket, error) {
	return s.bucketStore.Scan(ctx, func(k string) bool { return true })
}

func (s *GcsService) DeleteBucket(ctx context.Context, name string) error {
	if _, ok, _ := s.bucketStore.Get(ctx, name); !ok {
		return fmt.Errorf("NotFound: Bucket %s not found", name)
	}
	
	// TODO: Check if empty
	
	if err := s.bucketStore.Delete(ctx, name); err != nil {
		return err
	}
	
	return os.RemoveAll(filepath.Join(s.dataRoot, name))
}

// Objects

func (s *GcsService) InsertObject(ctx context.Context, bucketName, objectName string, content io.Reader, contentType string) (*model.ObjectMeta, error) {
	if _, err := s.GetBucket(ctx, bucketName); err != nil {
		return nil, err
	}

	// Save file
	path := filepath.Join(s.dataRoot, bucketName, objectName)
	os.MkdirAll(filepath.Dir(path), 0755)
	
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	hasher := md5.New()
	multi := io.MultiWriter(f, hasher)
	size, err := io.Copy(multi, content)
	if err != nil {
		return nil, err
	}

	md5Hash := hex.EncodeToString(hasher.Sum(nil))
	
	meta := &model.ObjectMeta{
		Kind:           "storage#object",
		Id:             fmt.Sprintf("%s/%s/1", bucketName, objectName),
		Name:           objectName,
		Bucket:         bucketName,
		SelfLink:       fmt.Sprintf("https://www.googleapis.com/storage/v1/b/%s/o/%s", bucketName, objectName),
		MediaLink:      fmt.Sprintf("https://www.googleapis.com/download/storage/v1/b/%s/o/%s?alt=media", bucketName, objectName),
		ContentType:    contentType,
		Size:           fmt.Sprintf("%d", size),
		TimeCreated:    time.Now(),
		Updated:        time.Now(),
		Generation:     "1",
		Metageneration: "1",
		Etag:           md5Hash,
		Md5Hash:        md5Hash,
		StorageClass:   "STANDARD",
	}

	key := fmt.Sprintf("%s:%s", bucketName, objectName)
	if err := s.objectStore.Put(ctx, key, meta); err != nil {
		return nil, err
	}

	return meta, nil
}

func (s *GcsService) GetObject(ctx context.Context, bucketName, objectName string) (*model.ObjectMeta, io.ReadCloser, error) {
	key := fmt.Sprintf("%s:%s", bucketName, objectName)
	meta, ok, err := s.objectStore.Get(ctx, key)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		return nil, nil, fmt.Errorf("NotFound: Object %s not found in bucket %s", objectName, bucketName)
	}

	path := filepath.Join(s.dataRoot, bucketName, objectName)
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	return meta, f, nil
}

func (s *GcsService) ListObjects(ctx context.Context, bucketName string) ([]*model.ObjectMeta, error) {
	return s.objectStore.Scan(ctx, func(k string) bool {
		meta, _, _ := s.objectStore.Get(ctx, k)
		return meta.Bucket == bucketName
	})
}

func (s *GcsService) DeleteObject(ctx context.Context, bucketName, objectName string) error {
	key := fmt.Sprintf("%s:%s", bucketName, objectName)
	if err := s.objectStore.Delete(ctx, key); err != nil {
		return err
	}
	
	path := filepath.Join(s.dataRoot, bucketName, objectName)
	return os.Remove(path)
}
