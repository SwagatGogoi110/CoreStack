package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/gcp/iam/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type IamGcpService struct {
	accountStore storage.Backend[string, *model.ServiceAccount]
}

func NewIamGcpService(factory *storage.Factory) (*IamGcpService, error) {
	accountStore, _ := storage.CreateAccountAware[*model.ServiceAccount](factory, "iam-gcp", "iam-accounts.json", "wal")

	return &IamGcpService{
		accountStore: accountStore,
	}, nil
}

func (s *IamGcpService) CreateServiceAccount(ctx context.Context, project, accountId, displayName string) (*model.ServiceAccount, error) {
	name := fmt.Sprintf("projects/%s/serviceAccounts/%s@%s.iam.gserviceaccount.com", project, accountId, project)
	
	if _, ok, _ := s.accountStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("AlreadyExists: Service account %s already exists", name)
	}

	sa := &model.ServiceAccount{
		Name:           name,
		ProjectId:      project,
		UniqueId:       uuid.New().String(),
		Email:          fmt.Sprintf("%s@%s.iam.gserviceaccount.com", accountId, project),
		DisplayName:    displayName,
		Oauth2ClientId: "123456789",
	}

	if err := s.accountStore.Put(ctx, name, sa); err != nil {
		return nil, err
	}

	return sa, nil
}

func (s *IamGcpService) ListServiceAccounts(ctx context.Context, project string) ([]*model.ServiceAccount, error) {
	prefix := fmt.Sprintf("projects/%s/", project)
	return s.accountStore.Scan(ctx, func(k string) bool {
		return strings.HasPrefix(k, prefix)
	})
}

func (s *IamGcpService) DeleteServiceAccount(ctx context.Context, name string) error {
	return s.accountStore.Delete(ctx, name)
}
