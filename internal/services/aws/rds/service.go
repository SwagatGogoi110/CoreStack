package rds

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/aws/rds/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type RdsService struct {
	instanceStore storage.Backend[string, *model.DbInstance]
	mu            sync.RWMutex
}

func NewRdsService(factory *storage.Factory) (*RdsService, error) {
	instanceStore, _ := storage.CreateAccountAware[*model.DbInstance](factory, "rds", "rds-instances.json", "wal")

	return &RdsService{
		instanceStore: instanceStore,
	}, nil
}

func (s *RdsService) CreateDBInstance(ctx context.Context, id, engine, class string, storageSize int) (*model.DbInstance, error) {
	if _, ok, _ := s.instanceStore.Get(ctx, id); ok {
		return nil, fmt.Errorf("DBInstanceAlreadyExists: DB instance already exists")
	}

	inst := &model.DbInstance{
		DBInstanceIdentifier: id,
		Engine:               engine,
		DBInstanceStatus:     "available",
		DBInstanceClass:      class,
		AllocatedStorage:     storageSize,
		InstanceCreateTime:   time.Now(),
		Endpoint: &model.Endpoint{
			Address: "localhost",
			Port:    5432, // Default to postgres for now
		},
	}

	if err := s.instanceStore.Put(ctx, id, inst); err != nil {
		return nil, err
	}

	log.Printf("Created RDS instance: %s", id)
	return inst, nil
}

func (s *RdsService) DescribeDBInstances(ctx context.Context, id string) ([]*model.DbInstance, error) {
	if id != "" {
		inst, ok, err := s.instanceStore.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("DBInstanceNotFound: DB instance not found")
		}
		return []*model.DbInstance{inst}, nil
	}
	return s.instanceStore.Scan(ctx, func(k string) bool { return true })
}
