package neptune

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/aws/neptune/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type NeptuneService struct {
	clusterStore storage.Backend[string, *model.NeptuneCluster]
	mu           sync.RWMutex
}

func NewNeptuneService(factory *storage.Factory) (*NeptuneService, error) {
	clusterStore, _ := storage.CreateAccountAware[*model.NeptuneCluster](factory, "neptune", "neptune-clusters.json", "wal")

	return &NeptuneService{
		clusterStore: clusterStore,
	}, nil
}

func (s *NeptuneService) CreateDBCluster(ctx context.Context, id, engineVersion string) (*model.NeptuneCluster, error) {
	if _, ok, _ := s.clusterStore.Get(ctx, id); ok {
		return nil, fmt.Errorf("DBClusterAlreadyExistsFault: Cluster already exists")
	}

	cluster := &model.NeptuneCluster{
		DBClusterIdentifier: id,
		DBClusterArn:        fmt.Sprintf("arn:aws:neptune:us-east-1:000000000000:cluster:%s", id),
		Status:              "available",
		Endpoint:            "localhost",
		Port:                8182,
		EngineVersion:       engineVersion,
		CreatedAt:           time.Now(),
	}

	if err := s.clusterStore.Put(ctx, id, cluster); err != nil {
		return nil, err
	}

	log.Printf("Created Neptune cluster: %s", id)
	return cluster, nil
}

func (s *NeptuneService) DescribeDBClusters(ctx context.Context, id string) ([]*model.NeptuneCluster, error) {
	if id != "" {
		c, ok, err := s.clusterStore.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("DBClusterNotFoundFault: Cluster not found")
		}
		return []*model.NeptuneCluster{c}, nil
	}
	return s.clusterStore.Scan(ctx, func(k string) bool { return true })
}
