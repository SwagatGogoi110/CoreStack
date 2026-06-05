package elasticache

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/elasticache/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type ElastiCacheService struct {
	groupStore storage.Backend[string, *model.ReplicationGroup]
	userStore  storage.Backend[string, *model.ElastiCacheUser]
	mu         sync.RWMutex
}

func NewElastiCacheService(factory *storage.Factory) (*ElastiCacheService, error) {
	groupStore, _ := storage.CreateAccountAware[*model.ReplicationGroup](factory, "elasticache", "ec-groups.json", "wal")
	userStore, _ := storage.CreateAccountAware[*model.ElastiCacheUser](factory, "elasticache", "ec-users.json", "wal")

	return &ElastiCacheService{
		groupStore: groupStore,
		userStore:  userStore,
	}, nil
}

func (s *ElastiCacheService) CreateReplicationGroup(ctx context.Context, id, description string) (*model.ReplicationGroup, error) {
	if _, ok, _ := s.groupStore.Get(ctx, id); ok {
		return nil, fmt.Errorf("ReplicationGroupAlreadyExistsFault: Replication group already exists")
	}

	group := &model.ReplicationGroup{
		ReplicationGroupID:     id,
		Description:            description,
		Status:                 "available",
		ReplicationGroupCreateTime: time.Now(),
		PrimaryEndpoint: &model.Endpoint{
			Address: "localhost",
			Port:    6379,
		},
	}

	if err := s.groupStore.Put(ctx, id, group); err != nil {
		return nil, err
	}

	log.Printf("Created ElastiCache replication group: %s", id)
	return group, nil
}

func (s *ElastiCacheService) DescribeReplicationGroups(ctx context.Context, id string) ([]*model.ReplicationGroup, error) {
	if id != "" {
		g, ok, err := s.groupStore.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("ReplicationGroupNotFoundFault: Replication group not found")
		}
		return []*model.ReplicationGroup{g}, nil
	}
	return s.groupStore.Scan(ctx, func(k string) bool { return true })
}
