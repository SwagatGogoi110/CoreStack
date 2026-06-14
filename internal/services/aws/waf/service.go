package waf

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/aws/waf/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type WafService struct {
	aclStore storage.Backend[string, *model.WebACL]
}

func NewWafService(factory *storage.Factory) (*WafService, error) {
	aclStore, _ := storage.CreateAccountAware[*model.WebACL](factory, "waf", "waf-acls.json", "wal")

	return &WafService{
		aclStore: aclStore,
	}, nil
}

func (s *WafService) CreateWebACL(ctx context.Context, name, metricName string, defaultAction string) (*model.WebACL, error) {
	id := uuid.New().String()
	acl := &model.WebACL{
		WebACLId:   id,
		Name:       name,
		MetricName: metricName,
		DefaultAction: &model.WafAction{
			Type: defaultAction,
		},
	}

	if err := s.aclStore.Put(ctx, id, acl); err != nil {
		return nil, err
	}
	return acl, nil
}

func (s *WafService) GetWebACL(ctx context.Context, id string) (*model.WebACL, error) {
	acl, ok, err := s.aclStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("WAFNonexistentItemException: WebACL %s not found", id)
	}
	return acl, nil
}

func (s *WafService) ListWebACLs(ctx context.Context) ([]*model.WebACL, error) {
	return s.aclStore.Scan(ctx, func(k string) bool { return true })
}
