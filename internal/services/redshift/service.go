package redshift

import (
	"context"
	"fmt"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/redshift/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type RedshiftService struct {
	clusterStore storage.Backend[string, *model.Cluster]
}

func NewRedshiftService(factory *storage.Factory) (*RedshiftService, error) {
	clusterStore, _ := storage.CreateAccountAware[*model.Cluster](factory, "redshift", "redshift-clusters.json", "wal")

	return &RedshiftService{
		clusterStore: clusterStore,
	}, nil
}

func (s *RedshiftService) CreateCluster(ctx context.Context, id, nodeType, masterUser, dbName string) (*model.Cluster, error) {
	if _, ok, _ := s.clusterStore.Get(ctx, id); ok {
		return nil, fmt.Errorf("ClusterAlreadyExists: Cluster already exists")
	}

	cluster := &model.Cluster{
		ClusterIdentifier: id,
		NodeType:          nodeType,
		ClusterStatus:     "available",
		MasterUsername:    masterUser,
		DBName:            dbName,
		ClusterCreateTime: time.Now(),
		Endpoint: &model.Endpoint{
			Address: fmt.Sprintf("%s.redshift.localhost", id),
			Port:    5439,
		},
	}

	if err := s.clusterStore.Put(ctx, id, cluster); err != nil {
		return nil, err
	}

	return cluster, nil
}

func (s *RedshiftService) DescribeClusters(ctx context.Context, id string) ([]*model.Cluster, error) {
	if id != "" {
		cluster, ok, err := s.clusterStore.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("ClusterNotFound: Cluster not found")
		}
		return []*model.Cluster{cluster}, nil
	}
	return s.clusterStore.Scan(ctx, func(k string) bool { return true })
}

func (s *RedshiftService) DeleteCluster(ctx context.Context, id string) error {
	_, ok, err := s.clusterStore.Get(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("ClusterNotFound: Cluster not found")
	}
	return s.clusterStore.Delete(ctx, id)
}
