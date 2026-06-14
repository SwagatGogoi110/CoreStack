package sagemaker

import (
	"context"
	"fmt"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/aws/sagemaker/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type SagemakerService struct {
	notebookStore storage.Backend[string, *model.NotebookInstance]
	modelStore    storage.Backend[string, *model.ModelSummary]
}

func NewSagemakerService(factory *storage.Factory) (*SagemakerService, error) {
	notebookStore, _ := storage.CreateAccountAware[*model.NotebookInstance](factory, "sagemaker", "notebooks.json", "wal")
	modelStore, _ := storage.CreateAccountAware[*model.ModelSummary](factory, "sagemaker", "models.json", "wal")

	return &SagemakerService{
		notebookStore: notebookStore,
		modelStore:    modelStore,
	}, nil
}

func (s *SagemakerService) CreateNotebookInstance(ctx context.Context, name, instanceType string) (*model.NotebookInstance, error) {
	if _, ok, _ := s.notebookStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("ResourceInUse: Notebook instance %s already exists", name)
	}

	nb := &model.NotebookInstance{
		NotebookInstanceName:   name,
		NotebookInstanceArn:    fmt.Sprintf("arn:aws:sagemaker:us-east-1:000000000000:notebook-instance/%s", name),
		NotebookInstanceStatus: "InService",
		InstanceType:           instanceType,
		CreationTime:           time.Now(),
		LastModifiedTime:       time.Now(),
	}

	s.notebookStore.Put(ctx, name, nb)
	return nb, nil
}

func (s *SagemakerService) ListNotebookInstances(ctx context.Context) ([]*model.NotebookInstance, error) {
	return s.notebookStore.Scan(ctx, func(k string) bool { return true })
}

func (s *SagemakerService) CreateModel(ctx context.Context, name string) (*model.ModelSummary, error) {
	m := &model.ModelSummary{
		ModelName:    name,
		ModelArn:     fmt.Sprintf("arn:aws:sagemaker:us-east-1:000000000000:model/%s", name),
		CreationTime: time.Now(),
	}
	s.modelStore.Put(ctx, name, m)
	return m, nil
}

func (s *SagemakerService) ListModels(ctx context.Context) ([]*model.ModelSummary, error) {
	return s.modelStore.Scan(ctx, func(k string) bool { return true })
}
