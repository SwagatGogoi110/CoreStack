package kms

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/kms/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type KmsService struct {
	keyStore storage.Backend[string, *model.KmsKey]
}

func NewKmsService(factory *storage.Factory) (*KmsService, error) {
	keyStore, _ := storage.CreateAccountAware[*model.KmsKey](factory, "kms", "kms-keys.json", "wal")

	return &KmsService{
		keyStore: keyStore,
	}, nil
}

func (s *KmsService) CreateKey(ctx context.Context, description string) (*model.KmsKey, error) {
	id := uuid.New().String()
	key := &model.KmsKey{
		KeyID:                 id,
		Arn:                   fmt.Sprintf("arn:aws:kms:us-east-1:000000000000:key/%s", id),
		Description:           description,
		Enabled:               true,
		KeyState:              "Enabled",
		KeyUsage:              "ENCRYPT_DECRYPT",
		CustomerMasterKeySpec: "SYMMETRIC_DEFAULT",
		CreationDate:          time.Now().Unix(),
	}

	if err := s.keyStore.Put(ctx, id, key); err != nil {
		return nil, err
	}

	log.Printf("Created KMS key: %s", id)
	return key, nil
}

func (s *KmsService) DescribeKey(ctx context.Context, id string) (*model.KmsKey, error) {
	key, ok, err := s.keyStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NotFoundException: Key not found")
	}
	return key, nil
}

func (s *KmsService) Encrypt(ctx context.Context, keyID string, plaintext []byte) ([]byte, error) {
	_, err := s.DescribeKey(ctx, keyID)
	if err != nil {
		return nil, err
	}
	// Mock encryption: base64 encode
	return []byte(base64.StdEncoding.EncodeToString(plaintext)), nil
}

func (s *KmsService) Decrypt(ctx context.Context, keyID string, ciphertext []byte) ([]byte, error) {
	// Mock decryption: base64 decode
	return base64.StdEncoding.DecodeString(string(ciphertext))
}
