package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/gcp/kafka/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type KafkaService struct {
	clusterStore storage.Backend[string, *model.Cluster]
}

func NewKafkaService(factory *storage.Factory) (*KafkaService, error) {
	clusterStore, _ := storage.CreateAccountAware[*model.Cluster](factory, "kafka", "kafka-clusters.json", "wal")

	return &KafkaService{
		clusterStore: clusterStore,
	}, nil
}

func (s *KafkaService) CreateCluster(ctx context.Context, name string, cluster *model.Cluster) (*model.Cluster, error) {
	if _, ok, _ := s.clusterStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("AlreadyExists: Cluster %s already exists", name)
	}

	cluster.Name = name
	cluster.State = "ACTIVE"
	cluster.CreateTime = time.Now()
	cluster.UpdateTime = time.Now()
	cluster.BootstrapAddress = "localhost:9092"

	if err := s.clusterStore.Put(ctx, name, cluster); err != nil {
		return nil, err
	}

	return cluster, nil
}

func (s *KafkaService) GetCluster(ctx context.Context, name string) (*model.Cluster, error) {
	c, ok, err := s.clusterStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NotFound: Cluster %s not found", name)
	}
	return c, nil
}

func (s *KafkaService) ListClusters(ctx context.Context, parent string) ([]*model.Cluster, error) {
	return s.clusterStore.Scan(ctx, func(k string) bool { return true })
}

func (s *KafkaService) DeleteCluster(ctx context.Context, name string) error {
	return s.clusterStore.Delete(ctx, name)
}
