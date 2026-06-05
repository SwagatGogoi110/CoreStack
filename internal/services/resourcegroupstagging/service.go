package resourcegroupstagging

import (
	"context"
	"sync"

	"github.com/hectorvent/cloudstack/internal/services/resourcegroupstagging/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type ResourceGroupsTaggingService struct {
	tagStore storage.Backend[string, *model.ResourceTagMapping]
	mu       sync.RWMutex
}

func NewResourceGroupsTaggingService(factory *storage.Factory) (*ResourceGroupsTaggingService, error) {
	tagStore, _ := storage.CreateAccountAware[*model.ResourceTagMapping](factory, "resourcegroupstagging", "resource-tags.json", "wal")

	return &ResourceGroupsTaggingService{
		tagStore: tagStore,
	}, nil
}

func (s *ResourceGroupsTaggingService) TagResources(ctx context.Context, arns []string, tags map[string]string) error {
	for _, arn := range arns {
		mapping, ok, _ := s.tagStore.Get(ctx, arn)
		if !ok {
			mapping = &model.ResourceTagMapping{
				ResourceArn: arn,
				Tags:        make(map[string]string),
			}
		}
		for k, v := range tags {
			mapping.Tags[k] = v
		}
		s.tagStore.Put(ctx, arn, mapping)
	}
	return nil
}

func (s *ResourceGroupsTaggingService) GetResources(ctx context.Context) ([]*model.ResourceTagMapping, error) {
	return s.tagStore.Scan(ctx, func(k string) bool { return true })
}
