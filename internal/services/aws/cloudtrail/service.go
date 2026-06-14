package cloudtrail

import (
	"context"
	"fmt"

	"github.com/hectorvent/cloudstack/internal/services/aws/cloudtrail/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CloudtrailService struct {
	trailStore storage.Backend[string, *model.Trail]
}

func NewCloudtrailService(factory *storage.Factory) (*CloudtrailService, error) {
	trailStore, _ := storage.CreateAccountAware[*model.Trail](factory, "cloudtrail", "cloudtrail-trails.json", "wal")

	return &CloudtrailService{
		trailStore: trailStore,
	}, nil
}

func (s *CloudtrailService) CreateTrail(ctx context.Context, name, bucketName string) (*model.Trail, error) {
	if _, ok, _ := s.trailStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("TrailAlreadyExistsException: Trail already exists")
	}

	trail := &model.Trail{
		Name:         name,
		S3BucketName: bucketName,
		HomeRegion:   "us-east-1",
		TrailARN:     fmt.Sprintf("arn:aws:cloudtrail:us-east-1:000000000000:trail/%s", name),
	}

	if err := s.trailStore.Put(ctx, name, trail); err != nil {
		return nil, err
	}

	return trail, nil
}

func (s *CloudtrailService) DescribeTrails(ctx context.Context, names []string) ([]*model.Trail, error) {
	if len(names) > 0 {
		var trails []*model.Trail
		for _, name := range names {
			trail, ok, err := s.trailStore.Get(ctx, name)
			if err == nil && ok {
				trails = append(trails, trail)
			}
		}
		return trails, nil
	}
	return s.trailStore.Scan(ctx, func(k string) bool { return true })
}

func (s *CloudtrailService) DeleteTrail(ctx context.Context, name string) error {
	return s.trailStore.Delete(ctx, name)
}
