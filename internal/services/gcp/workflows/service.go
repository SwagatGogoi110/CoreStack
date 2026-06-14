package workflows

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/gcp/workflows/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type WorkflowsService struct {
	workflowStore  storage.Backend[string, *model.Workflow]
	executionStore storage.Backend[string, *model.Execution]
}

func NewWorkflowsService(factory *storage.Factory) (*WorkflowsService, error) {
	workflowStore, _ := storage.CreateAccountAware[*model.Workflow](factory, "workflows", "wf-workflows.json", "wal")
	executionStore, _ := storage.CreateAccountAware[*model.Execution](factory, "workflows", "wf-executions.json", "wal")

	return &WorkflowsService{
		workflowStore:  workflowStore,
		executionStore: executionStore,
	}, nil
}

// Workflows

func (s *WorkflowsService) CreateWorkflow(ctx context.Context, name string, wf *model.Workflow) (*model.Workflow, error) {
	wf.Name = name
	wf.State = "ACTIVE"
	wf.CreateTime = time.Now()
	wf.UpdateTime = time.Now()

	if err := s.workflowStore.Put(ctx, name, wf); err != nil {
		return nil, err
	}
	return wf, nil
}

func (s *WorkflowsService) ListWorkflows(ctx context.Context, parent string) ([]*model.Workflow, error) {
	return s.workflowStore.Scan(ctx, func(k string) bool { return true })
}

// Executions

func (s *WorkflowsService) CreateExecution(ctx context.Context, workflowName string, execution *model.Execution) (*model.Execution, error) {
	id := uuid.New().String()
	execution.Name = workflowName + "/executions/" + id
	execution.State = "SUCCEEDED"
	execution.StartTime = time.Now()
	execution.EndTime = time.Now()
	execution.Result = "{\"status\": \"ok\"}"

	if err := s.executionStore.Put(ctx, execution.Name, execution); err != nil {
		return nil, err
	}
	return execution, nil
}

func (s *WorkflowsService) ListExecutions(ctx context.Context, workflowName string) ([]*model.Execution, error) {
	return s.executionStore.Scan(ctx, func(k string) bool { return true })
}
