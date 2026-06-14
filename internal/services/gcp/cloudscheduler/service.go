package cloudscheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/gcp/cloudscheduler/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CloudSchedulerService struct {
	jobStore storage.Backend[string, *model.Job]
}

func NewCloudSchedulerService(factory *storage.Factory) (*CloudSchedulerService, error) {
	jobStore, _ := storage.CreateAccountAware[*model.Job](factory, "cloudscheduler", "cs-jobs.json", "wal")

	return &CloudSchedulerService{
		jobStore: jobStore,
	}, nil
}

// Jobs

func (s *CloudSchedulerService) CreateJob(ctx context.Context, name string, job *model.Job) (*model.Job, error) {
	job.Name = name
	job.State = "ENABLED"
	if err := s.jobStore.Put(ctx, name, job); err != nil {
		return nil, err
	}
	return job, nil
}

func (s *CloudSchedulerService) ListJobs(ctx context.Context, parent string) ([]*model.Job, error) {
	return s.jobStore.Scan(ctx, func(k string) bool { return true })
}

func (s *CloudSchedulerService) RunJob(ctx context.Context, name string) (*model.Job, error) {
	job, ok, err := s.jobStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NotFound: Job %s not found", name)
	}
	job.LastAttemptTime = time.Now()
	job.Status = &model.Status{Code: 0, Message: "Success"}
	s.jobStore.Put(ctx, name, job)
	return job, nil
}

func (s *CloudSchedulerService) DeleteJob(ctx context.Context, name string) error {
	return s.jobStore.Delete(ctx, name)
}
