package operations

import (
	"context"
	"fmt"

	"github.com/hectorvent/cloudstack/internal/services/gcp/operations/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type OperationsService struct {
	opStore storage.Backend[string, *model.Operation]
}

func NewOperationsService(factory *storage.Factory) (*OperationsService, error) {
	opStore, _ := storage.CreateAccountAware[*model.Operation](factory, "operations", "sm-operations.json", "wal")

	return &OperationsService{
		opStore: opStore,
	}, nil
}

func (s *OperationsService) GetOperation(ctx context.Context, name string) (*model.Operation, error) {
	op, ok, err := s.opStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NotFound: Operation %s not found", name)
	}
	return op, nil
}

func (s *OperationsService) ListOperations(ctx context.Context, name string) ([]*model.Operation, error) {
	return s.opStore.Scan(ctx, func(k string) bool { return true })
}

func (s *OperationsService) DeleteOperation(ctx context.Context, name string) error {
	return s.opStore.Delete(ctx, name)
}
