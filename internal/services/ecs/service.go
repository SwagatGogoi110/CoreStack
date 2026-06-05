package ecs

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/hectorvent/cloudstack/internal/services/ecs/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type EcsService struct {
	clusterStore storage.Backend[string, *model.EcsCluster]
	mu           sync.RWMutex
}

func NewEcsService(factory *storage.Factory) (*EcsService, error) {
	clusterStore, _ := storage.CreateAccountAware[*model.EcsCluster](factory, "ecs", "ecs-clusters.json", "wal")

	return &EcsService{
		clusterStore: clusterStore,
	}, nil
}

func (s *EcsService) CreateCluster(ctx context.Context, name string) (*model.EcsCluster, error) {
	if _, ok, _ := s.clusterStore.Get(ctx, name); ok {
		// Existing
	}

	arn := fmt.Sprintf("arn:aws:ecs:us-east-1:000000000000:cluster/%s", name)
	cluster := &model.EcsCluster{
		ClusterName: name,
		ClusterArn:  arn,
		Status:      "ACTIVE",
	}

	if err := s.clusterStore.Put(ctx, name, cluster); err != nil {
		return nil, err
	}

	log.Printf("Created ECS cluster: %s", name)
	return cluster, nil
}

func (s *EcsService) DescribeClusters(ctx context.Context, names []string) ([]*model.EcsCluster, error) {
	return s.clusterStore.Scan(ctx, func(k string) bool {
		if len(names) == 0 {
			return true
		}
		for _, n := range names {
			if n == k {
				return true
			}
		}
		return false
	})
}
