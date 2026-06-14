package container

import (
	"context"
	"fmt"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/gcp/container/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type GkeService struct {
	clusterStore  storage.Backend[string, *model.Cluster]
	nodePoolStore storage.Backend[string, *model.NodePool]
}

func NewGkeService(factory *storage.Factory) (*GkeService, error) {
	clusterStore, _ := storage.CreateAccountAware[*model.Cluster](factory, "container", "gke-clusters.json", "wal")
	nodePoolStore, _ := storage.CreateAccountAware[*model.NodePool](factory, "container", "gke-node-pools.json", "wal")

	return &GkeService{
		clusterStore:  clusterStore,
		nodePoolStore: nodePoolStore,
	}, nil
}

// Clusters

func (s *GkeService) CreateCluster(ctx context.Context, name string, cluster *model.Cluster) (*model.Cluster, error) {
	if _, ok, _ := s.clusterStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("AlreadyExists: Cluster %s already exists", name)
	}

	cluster.Name = name
	cluster.Status = "RUNNING"
	cluster.CreateTime = time.Now()
	cluster.Endpoint = "127.0.0.1"

	if err := s.clusterStore.Put(ctx, name, cluster); err != nil {
		return nil, err
	}
	return cluster, nil
}

func (s *GkeService) ListClusters(ctx context.Context, parent string) ([]*model.Cluster, error) {
	return s.clusterStore.Scan(ctx, func(k string) bool { return true })
}

// Node Pools

func (s *GkeService) CreateNodePool(ctx context.Context, parent, name string, pool *model.NodePool) (*model.NodePool, error) {
	pool.Name = name
	pool.Status = "RUNNING"
	if err := s.nodePoolStore.Put(ctx, parent+"/"+name, pool); err != nil {
		return nil, err
	}
	return pool, nil
}

func (s *GkeService) ListNodePools(ctx context.Context, parent string) ([]*model.NodePool, error) {
	return s.nodePoolStore.Scan(ctx, func(k string) bool { return true })
}
