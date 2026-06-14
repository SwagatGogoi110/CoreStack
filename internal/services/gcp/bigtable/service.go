package bigtable

import (
	"context"
	"fmt"

	"github.com/hectorvent/cloudstack/internal/services/gcp/bigtable/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type BigtableService struct {
	instanceStore storage.Backend[string, *model.Instance]
	tableStore    storage.Backend[string, *model.Table]
}

func NewBigtableService(factory *storage.Factory) (*BigtableService, error) {
	instanceStore, _ := storage.CreateAccountAware[*model.Instance](factory, "bigtable", "bt-instances.json", "wal")
	tableStore, _ := storage.CreateAccountAware[*model.Table](factory, "bigtable", "bt-tables.json", "wal")

	return &BigtableService{
		instanceStore: instanceStore,
		tableStore:    tableStore,
	}, nil
}

// Instances

func (s *BigtableService) CreateInstance(ctx context.Context, name string, inst *model.Instance) (*model.Instance, error) {
	if _, ok, _ := s.instanceStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("AlreadyExists: Instance %s already exists", name)
	}

	inst.Name = name
	inst.State = "READY"

	if err := s.instanceStore.Put(ctx, name, inst); err != nil {
		return nil, err
	}
	return inst, nil
}

func (s *BigtableService) ListInstances(ctx context.Context, project string) ([]*model.Instance, error) {
	return s.instanceStore.Scan(ctx, func(k string) bool { return true })
}

// Tables

func (s *BigtableService) CreateTable(ctx context.Context, name string, table *model.Table) (*model.Table, error) {
	table.Name = name
	if err := s.tableStore.Put(ctx, name, table); err != nil {
		return nil, err
	}
	return table, nil
}

func (s *BigtableService) ListTables(ctx context.Context, instanceName string) ([]*model.Table, error) {
	return s.tableStore.Scan(ctx, func(k string) bool { return true })
}
