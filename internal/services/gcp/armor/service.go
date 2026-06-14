package armor

import (
	"context"
	"fmt"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/gcp/armor/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CloudArmorService struct {
	policyStore storage.Backend[string, *model.SecurityPolicy]
}

func NewCloudArmorService(factory *storage.Factory) (*CloudArmorService, error) {
	policyStore, _ := storage.CreateAccountAware[*model.SecurityPolicy](factory, "armor", "armor-policies.json", "wal")

	return &CloudArmorService{
		policyStore: policyStore,
	}, nil
}

// Policies

func (s *CloudArmorService) CreatePolicy(ctx context.Context, name string, policy *model.SecurityPolicy) (*model.SecurityPolicy, error) {
	if _, ok, _ := s.policyStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("AlreadyExists: Policy %s already exists", name)
	}

	policy.Name = name
	policy.CreateTime = time.Now()
	policy.SelfLink = fmt.Sprintf("https://www.googleapis.com/compute/v1/global/securityPolicies/%s", name)

	if err := s.policyStore.Put(ctx, name, policy); err != nil {
		return nil, err
	}
	return policy, nil
}

func (s *CloudArmorService) ListPolicies(ctx context.Context) ([]*model.SecurityPolicy, error) {
	return s.policyStore.Scan(ctx, func(k string) bool { return true })
}

func (s *CloudArmorService) DeletePolicy(ctx context.Context, name string) error {
	return s.policyStore.Delete(ctx, name)
}
