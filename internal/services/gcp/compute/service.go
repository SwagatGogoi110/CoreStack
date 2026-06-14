package compute

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/gcp/compute/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type ComputeEngineService struct {
	instanceStore storage.Backend[string, *model.Instance]
}

func NewComputeEngineService(factory *storage.Factory) (*ComputeEngineService, error) {
	instanceStore, _ := storage.CreateAccountAware[*model.Instance](factory, "compute", "ce-instances.json", "wal")

	return &ComputeEngineService{
		instanceStore: instanceStore,
	}, nil
}

// Instances

func (s *ComputeEngineService) CreateInstance(ctx context.Context, project, zone, name string, inst *model.Instance) (*model.Instance, error) {
	key := fmt.Sprintf("%s:%s:%s", project, zone, name)
	if _, ok, _ := s.instanceStore.Get(ctx, key); ok {
		return nil, fmt.Errorf("AlreadyExists: Instance %s already exists", name)
	}

	inst.Name = name
	inst.Id = uuid.New().String()
	inst.Kind = "compute#instance"
	inst.Status = "RUNNING"
	inst.CreationTimestamp = time.Now()
	inst.Zone = fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s", project, zone)
	inst.SelfLink = fmt.Sprintf("%s/instances/%s", inst.Zone, name)

	if err := s.instanceStore.Put(ctx, key, inst); err != nil {
		return nil, err
	}
	return inst, nil
}

func (s *ComputeEngineService) ListInstances(ctx context.Context, project, zone string) ([]*model.Instance, error) {
	prefix := fmt.Sprintf("%s:%s", project, zone)
	return s.instanceStore.Scan(ctx, func(k string) bool {
		return strings.HasPrefix(k, prefix)
	})
}

func (s *ComputeEngineService) DeleteInstance(ctx context.Context, project, zone, name string) error {
	key := fmt.Sprintf("%s:%s:%s", project, zone, name)
	return s.instanceStore.Delete(ctx, key)
}
