package synthetics

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/synthetics/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type SyntheticsService struct {
	canaryStore storage.Backend[string, *model.Canary]
}

func NewSyntheticsService(factory *storage.Factory) (*SyntheticsService, error) {
	canaryStore, _ := storage.CreateAccountAware[*model.Canary](factory, "synthetics", "synthetics-canaries.json", "wal")

	return &SyntheticsService{
		canaryStore: canaryStore,
	}, nil
}

func (s *SyntheticsService) CreateCanary(ctx context.Context, name, artifactS3Location, executionRoleArn string) (*model.Canary, error) {
	if _, ok, _ := s.canaryStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("ConflictException: Canary %s already exists", name)
	}

	canary := &model.Canary{
		Id:                 uuid.New().String(),
		Name:               name,
		ArtifactS3Location: artifactS3Location,
		ExecutionRoleArn:    executionRoleArn,
		Status: &model.CanaryStatus{
			State: "READY",
		},
	}

	s.canaryStore.Put(ctx, name, canary)
	return canary, nil
}

func (s *SyntheticsService) GetCanary(ctx context.Context, name string) (*model.Canary, error) {
	canary, ok, err := s.canaryStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("ResourceNotFoundException: Canary %s not found", name)
	}
	return canary, nil
}

func (s *SyntheticsService) DescribeCanaries(ctx context.Context) ([]*model.Canary, error) {
	return s.canaryStore.Scan(ctx, func(k string) bool { return true })
}

func (s *SyntheticsService) DeleteCanary(ctx context.Context, name string) error {
	return s.canaryStore.Delete(ctx, name)
}
