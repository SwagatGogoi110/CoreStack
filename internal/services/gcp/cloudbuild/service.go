package cloudbuild

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/gcp/cloudbuild/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CloudBuildService struct {
	buildStore   storage.Backend[string, *model.Build]
	triggerStore storage.Backend[string, *model.BuildTrigger]
}

func NewCloudBuildService(factory *storage.Factory) (*CloudBuildService, error) {
	buildStore, _ := storage.CreateAccountAware[*model.Build](factory, "cloudbuild", "cb-builds.json", "wal")
	triggerStore, _ := storage.CreateAccountAware[*model.BuildTrigger](factory, "cloudbuild", "cb-triggers.json", "wal")

	return &CloudBuildService{
		buildStore:   buildStore,
		triggerStore: triggerStore,
	}, nil
}

// Builds

func (s *CloudBuildService) CreateBuild(ctx context.Context, project string, build *model.Build) (*model.Build, error) {
	id := uuid.New().String()
	build.Id = id
	build.ProjectId = project
	build.Status = "SUCCESS"
	build.CreateTime = time.Now()
	build.StartTime = time.Now()
	build.FinishTime = time.Now()
	build.LogUrl = fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds/%s?project=%s", id, project)

	if err := s.buildStore.Put(ctx, id, build); err != nil {
		return nil, err
	}
	return build, nil
}

func (s *CloudBuildService) ListBuilds(ctx context.Context, project string) ([]*model.Build, error) {
	return s.buildStore.Scan(ctx, func(k string) bool { return true })
}

// Triggers

func (s *CloudBuildService) CreateTrigger(ctx context.Context, project string, trigger *model.BuildTrigger) (*model.BuildTrigger, error) {
	id := uuid.New().String()
	trigger.Id = id
	trigger.CreateTime = time.Now()
	if err := s.triggerStore.Put(ctx, id, trigger); err != nil {
		return nil, err
	}
	return trigger, nil
}

func (s *CloudBuildService) ListTriggers(ctx context.Context, project string) ([]*model.BuildTrigger, error) {
	return s.triggerStore.Scan(ctx, func(k string) bool { return true })
}
