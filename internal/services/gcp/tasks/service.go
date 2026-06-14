package tasks

import (
	"context"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/gcp/tasks/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type TasksService struct {
	queueStore storage.Backend[string, *model.Queue]
	taskStore  storage.Backend[string, *model.Task]
}

func NewTasksService(factory *storage.Factory) (*TasksService, error) {
	queueStore, _ := storage.CreateAccountAware[*model.Queue](factory, "tasks", "tasks-queues.json", "wal")
	taskStore, _ := storage.CreateAccountAware[*model.Task](factory, "tasks", "tasks-tasks.json", "wal")

	return &TasksService{
		queueStore: queueStore,
		taskStore:  taskStore,
	}, nil
}

func (s *TasksService) CreateQueue(ctx context.Context, name string, queue *model.Queue) (*model.Queue, error) {
	queue.Name = name
	queue.State = "RUNNING"
	if err := s.queueStore.Put(ctx, name, queue); err != nil {
		return nil, err
	}
	return queue, nil
}

func (s *TasksService) ListQueues(ctx context.Context, parent string) ([]*model.Queue, error) {
	return s.queueStore.Scan(ctx, func(k string) bool { return true })
}

func (s *TasksService) CreateTask(ctx context.Context, parent string, task *model.Task) (*model.Task, error) {
	task.CreateTime = time.Now()
	if task.ScheduleTime.IsZero() {
		task.ScheduleTime = time.Now()
	}
	
	if err := s.taskStore.Put(ctx, task.Name, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TasksService) ListTasks(ctx context.Context, parent string) ([]*model.Task, error) {
	return s.taskStore.Scan(ctx, func(k string) bool { return true })
}

func (s *TasksService) DeleteQueue(ctx context.Context, name string) error {
	return s.queueStore.Delete(ctx, name)
}
