package emr

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/emr/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type EmrService struct {
	clusterStore storage.Backend[string, *model.Cluster]
}

func NewEmrService(factory *storage.Factory) (*EmrService, error) {
	clusterStore, _ := storage.CreateAccountAware[*model.Cluster](factory, "emr", "emr-clusters.json", "wal")

	return &EmrService{
		clusterStore: clusterStore,
	}, nil
}

func (s *EmrService) RunJobFlow(ctx context.Context, name, releaseLabel string) (*model.Cluster, error) {
	id := "j-" + uuid.New().String()[:10]
	
	cluster := &model.Cluster{
		Id:           id,
		Name:         name,
		ReleaseLabel: releaseLabel,
		Status: &model.ClusterStatus{
			State: "WAITING",
			Timeline: &model.Timeline{
				CreationDateTime: time.Now(),
				ReadyDateTime:    time.Now().Add(time.Minute),
			},
		},
	}

	if err := s.clusterStore.Put(ctx, id, cluster); err != nil {
		return nil, err
	}

	return cluster, nil
}

func (s *EmrService) DescribeCluster(ctx context.Context, id string) (*model.Cluster, error) {
	cluster, ok, err := s.clusterStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("InvalidRequestException: Cluster with id %s not found", id)
	}
	return cluster, nil
}

func (s *EmrService) ListClusters(ctx context.Context) ([]*model.Cluster, error) {
	return s.clusterStore.Scan(ctx, func(k string) bool { return true })
}

func (s *EmrService) TerminateJobFlows(ctx context.Context, ids []string) error {
	for _, id := range ids {
		cluster, ok, err := s.clusterStore.Get(ctx, id)
		if err == nil && ok {
			cluster.Status.State = "TERMINATED"
			s.clusterStore.Put(ctx, id, cluster)
		}
	}
	return nil
}
