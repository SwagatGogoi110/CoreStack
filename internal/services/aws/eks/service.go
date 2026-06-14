package eks

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/aws/eks/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type EksService struct {
	clusterStore storage.Backend[string, *model.Cluster]
}

func NewEksService(factory *storage.Factory) (*EksService, error) {
	clusterStore, _ := storage.CreateAccountAware[*model.Cluster](factory, "eks", "eks-clusters.json", "wal")

	return &EksService{
		clusterStore: clusterStore,
	}, nil
}

func (s *EksService) CreateCluster(ctx context.Context, name, version string) (*model.Cluster, error) {
	if _, ok, _ := s.clusterStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("ResourceInUseException: Cluster already exists")
	}

	cluster := &model.Cluster{
		Name:      name,
		Arn:       fmt.Sprintf("arn:aws:eks:us-east-1:000000000000:cluster/%s", name),
		CreatedAt: time.Now(),
		Version:   version,
		Status:    "ACTIVE", // Auto-succeed in mock
		Endpoint:  fmt.Sprintf("https://%s.eks.localhost", name),
	}

	if err := s.clusterStore.Put(ctx, name, cluster); err != nil {
		return nil, err
	}

	log.Printf("Created EKS cluster: %s", name)
	return cluster, nil
}

func (s *EksService) DescribeCluster(ctx context.Context, name string) (*model.Cluster, error) {
	cluster, ok, err := s.clusterStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("ResourceNotFoundException: Cluster not found")
	}
	return cluster, nil
}

func (s *EksService) ListClusters(ctx context.Context) ([]string, error) {
	clusters, err := s.clusterStore.Scan(ctx, func(k string) bool { return true })
	if err != nil {
		return nil, err
	}
	var names []string
	for _, c := range clusters {
		names = append(names, c.Name)
	}
	return names, nil
}
