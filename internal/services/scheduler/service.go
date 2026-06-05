package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/scheduler/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type SchedulerService struct {
	groupStore    storage.Backend[string, *model.ScheduleGroup]
	scheduleStore storage.Backend[string, *model.Schedule]
	mu            sync.RWMutex
}

func NewSchedulerService(factory *storage.Factory) (*SchedulerService, error) {
	groupStore, _ := storage.CreateAccountAware[*model.ScheduleGroup](factory, "scheduler", "scheduler-groups.json", "wal")
	scheduleStore, _ := storage.CreateAccountAware[*model.Schedule](factory, "scheduler", "scheduler-schedules.json", "wal")

	return &SchedulerService{
		groupStore:    groupStore,
		scheduleStore: scheduleStore,
	}, nil
}

func (s *SchedulerService) CreateSchedule(ctx context.Context, name, group, expression string) (*model.Schedule, error) {
	if group == "" {
		group = "default"
	}
	key := fmt.Sprintf("%s:%s", group, name)
	if _, ok, _ := s.scheduleStore.Get(ctx, key); ok {
		return nil, fmt.Errorf("ConflictException: Schedule already exists")
	}

	schedule := &model.Schedule{
		Name:               name,
		GroupName:          group,
		Arn:                fmt.Sprintf("arn:aws:scheduler:us-east-1:000000000000:schedule/%s/%s", group, name),
		State:              "ENABLED",
		ScheduleExpression: expression,
		CreationDate:       time.Now(),
		LastModificationDate: time.Now(),
	}

	if err := s.scheduleStore.Put(ctx, key, schedule); err != nil {
		return nil, err
	}

	log.Printf("Created Schedule: %s in group %s", name, group)
	return schedule, nil
}

func (s *SchedulerService) GetSchedule(ctx context.Context, name, group string) (*model.Schedule, error) {
	if group == "" {
		group = "default"
	}
	key := fmt.Sprintf("%s:%s", group, name)
	schedule, ok, err := s.scheduleStore.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("ResourceNotFoundException: Schedule not found")
	}
	return schedule, nil
}

func (s *SchedulerService) ListSchedules(ctx context.Context) ([]*model.Schedule, error) {
	return s.scheduleStore.Scan(ctx, func(k string) bool { return true })
}
