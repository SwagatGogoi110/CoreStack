package secretmanager

import (
	"context"
	"fmt"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/gcp/secretmanager/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type SecretManagerGcpService struct {
	secretStore  storage.Backend[string, *model.Secret]
	versionStore storage.Backend[string, *model.SecretVersion]
	payloadStore storage.Backend[string, []byte]
}

func NewSecretManagerGcpService(factory *storage.Factory) (*SecretManagerGcpService, error) {
	secretStore, _ := storage.CreateAccountAware[*model.Secret](factory, "secretmanager", "sm-secrets.json", "wal")
	versionStore, _ := storage.CreateAccountAware[*model.SecretVersion](factory, "secretmanager", "sm-versions.json", "wal")
	payloadStore, _ := storage.CreateAccountAware[[]byte](factory, "secretmanager", "sm-payloads.json", "wal")

	return &SecretManagerGcpService{
		secretStore:  secretStore,
		versionStore: versionStore,
		payloadStore: payloadStore,
	}, nil
}

func (s *SecretManagerGcpService) CreateSecret(ctx context.Context, name string) (*model.Secret, error) {
	secret := &model.Secret{
		Name:       name,
		CreateTime: time.Now(),
		Replication: &model.Replication{
			Automatic: &model.Automatic{},
		},
	}

	if err := s.secretStore.Put(ctx, name, secret); err != nil {
		return nil, err
	}

	return secret, nil
}

func (s *SecretManagerGcpService) GetSecret(ctx context.Context, name string) (*model.Secret, error) {
	secret, ok, err := s.secretStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NotFound: Secret %s not found", name)
	}
	return secret, nil
}

func (s *SecretManagerGcpService) ListSecrets(ctx context.Context) ([]*model.Secret, error) {
	return s.secretStore.Scan(ctx, func(k string) bool { return true })
}

func (s *SecretManagerGcpService) AddSecretVersion(ctx context.Context, secretName string, data []byte) (*model.SecretVersion, error) {
	if _, err := s.GetSecret(ctx, secretName); err != nil {
		return nil, err
	}

	// Simple versioning
	versions, _ := s.versionStore.Scan(ctx, func(k string) bool {
		return len(k) > len(secretName) && k[:len(secretName)] == secretName
	})

	versionId := fmt.Sprintf("%d", len(versions)+1)
	name := secretName + "/versions/" + versionId

	version := &model.SecretVersion{
		Name:       name,
		CreateTime: time.Now(),
		State:      "ENABLED",
	}

	if err := s.versionStore.Put(ctx, name, version); err != nil {
		return nil, err
	}
	
	s.payloadStore.Put(ctx, name, data)

	return version, nil
}

func (s *SecretManagerGcpService) AccessSecretVersion(ctx context.Context, name string) (*model.SecretVersion, []byte, error) {
	version, ok, err := s.versionStore.Get(ctx, name)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		return nil, nil, fmt.Errorf("NotFound: Secret version %s not found", name)
	}

	payload, _, _ := s.payloadStore.Get(ctx, name)
	return version, payload, nil
}

func (s *SecretManagerGcpService) DeleteSecret(ctx context.Context, name string) error {
	return s.secretStore.Delete(ctx, name)
}
