package msk

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/aws/msk/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type MskService struct {
	clusterStore storage.Backend[string, *model.MskCluster]
}

func NewMskService(factory *storage.Factory) (*MskService, error) {
	clusterStore, _ := storage.CreateAccountAware[*model.MskCluster](factory, "msk", "msk-clusters.json", "wal")

	return &MskService{
		clusterStore: clusterStore,
	}, nil
}

func (s *MskService) CreateCluster(ctx context.Context, name string) (*model.MskCluster, error) {
	// Check if exists
	clusters, _ := s.clusterStore.Scan(ctx, func(k string) bool { return true })
	for _, c := range clusters {
		if c.ClusterName == name {
			return nil, fmt.Errorf("ConflictException: Cluster already exists")
		}
	}

	arn := fmt.Sprintf("arn:aws:kafka:us-east-1:000000000000:cluster/%s/%s", name, uuid.New().String())
	cluster := &model.MskCluster{
		ClusterArn:         arn,
		ClusterName:        name,
		State:              "ACTIVE",
		CreationTime:       time.Now(),
		CurrentVersion:     "3.6.0",
		NumberOfBrokerNodes: 1,
		ZookeeperConnectString: "localhost:2181",
	}

	if err := s.clusterStore.Put(ctx, arn, cluster); err != nil {
		return nil, err
	}

	log.Printf("Created MSK cluster: %s", name)
	return cluster, nil
}

func (s *MskService) DescribeCluster(ctx context.Context, arn string) (*model.MskCluster, error) {
	cluster, ok, err := s.clusterStore.Get(ctx, arn)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NotFoundException: Cluster not found")
	}
	return cluster, nil
}

func (s *MskService) ListClusters(ctx context.Context) ([]*model.MskCluster, error) {
	return s.clusterStore.Scan(ctx, func(k string) bool { return true })
}
