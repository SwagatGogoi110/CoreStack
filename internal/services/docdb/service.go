package docdb

import (
	"context"
	"fmt"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/docdb/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type DocdbService struct {
	clusterStore storage.Backend[string, *model.DBCluster]
}

func NewDocdbService(factory *storage.Factory) (*DocdbService, error) {
	clusterStore, _ := storage.CreateAccountAware[*model.DBCluster](factory, "docdb", "docdb-clusters.json", "wal")

	return &DocdbService{
		clusterStore: clusterStore,
	}, nil
}

func (s *DocdbService) CreateDBCluster(ctx context.Context, id, engine string) (*model.DBCluster, error) {
	if _, ok, _ := s.clusterStore.Get(ctx, id); ok {
		return nil, fmt.Errorf("DBClusterAlreadyExists: DB cluster already exists")
	}

	cluster := &model.DBCluster{
		DBClusterIdentifier: id,
		Engine:               engine,
		Status:               "available",
		Endpoint:             fmt.Sprintf("%s.docdb.localhost", id),
		Port:                 27017,
		ClusterCreateTime:   time.Now(),
	}

	if err := s.clusterStore.Put(ctx, id, cluster); err != nil {
		return nil, err
	}

	return cluster, nil
}

func (s *DocdbService) DescribeDBClusters(ctx context.Context, id string) ([]*model.DBCluster, error) {
	if id != "" {
		cluster, ok, err := s.clusterStore.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("DBClusterNotFound: DB cluster not found")
		}
		return []*model.DBCluster{cluster}, nil
	}
	return s.clusterStore.Scan(ctx, func(k string) bool { return true })
}

func (s *DocdbService) DeleteDBCluster(ctx context.Context, id string) error {
	return s.clusterStore.Delete(ctx, id)
}
