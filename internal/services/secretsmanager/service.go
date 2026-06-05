package secretsmanager

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/secretsmanager/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type SecretsManagerService struct {
	secretStore storage.Backend[string, *model.Secret]
}

func NewSecretsManagerService(factory *storage.Factory) (*SecretsManagerService, error) {
	secretStore, _ := storage.CreateAccountAware[*model.Secret](factory, "secretsmanager", "sm-secrets.json", "wal")

	return &SecretsManagerService{
		secretStore: secretStore,
	}, nil
}

func (s *SecretsManagerService) CreateSecret(ctx context.Context, name, value string) (*model.Secret, error) {
	if _, ok, _ := s.secretStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("ResourceExistsException: Secret already exists")
	}

	versionId := uuid.New().String()
	arn := fmt.Sprintf("arn:aws:secretsmanager:us-east-1:000000000000:secret:%s", name)

	secret := &model.Secret{
		Name:            name,
		Arn:             arn,
		CreatedDate:     time.Now(),
		LastChangedDate: time.Now(),
		Versions: map[string]*model.SecretVersion{
			versionId: {
				VersionId:     versionId,
				SecretString:  value,
				VersionStages: []string{"AWSCURRENT"},
				CreatedDate:   time.Now(),
			},
		},
		CurrentVersionId: versionId,
	}

	if err := s.secretStore.Put(ctx, name, secret); err != nil {
		return nil, err
	}

	log.Printf("Created Secret: %s", name)
	return secret, nil
}

func (s *SecretsManagerService) GetSecretValue(ctx context.Context, id string) (*model.SecretVersion, error) {
	secret, ok, err := s.secretStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("ResourceNotFoundException: Secret not found")
	}

	version := secret.Versions[secret.CurrentVersionId]
	return version, nil
}

func (s *SecretsManagerService) ListSecrets(ctx context.Context) ([]*model.Secret, error) {
	return s.secretStore.Scan(ctx, func(k string) bool { return true })
}
