package cloudrun

import (
	"context"
	"fmt"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/gcp/cloudrun/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CloudRunService struct {
	serviceStore storage.Backend[string, *model.Service]
}

func NewCloudRunService(factory *storage.Factory) (*CloudRunService, error) {
	serviceStore, _ := storage.CreateAccountAware[*model.Service](factory, "cloudrun", "cr-services.json", "wal")

	return &CloudRunService{
		serviceStore: serviceStore,
	}, nil
}

func (s *CloudRunService) CreateService(ctx context.Context, name string, svc *model.Service) (*model.Service, error) {
	if _, ok, _ := s.serviceStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("AlreadyExists: Service %s already exists", name)
	}

	svc.Name = name
	svc.Uri = fmt.Sprintf("https://%s.run.app", name)
	svc.CreateTime = time.Now()
	svc.UpdateTime = time.Now()

	if err := s.serviceStore.Put(ctx, name, svc); err != nil {
		return nil, err
	}

	return svc, nil
}

func (s *CloudRunService) GetService(ctx context.Context, name string) (*model.Service, error) {
	svc, ok, err := s.serviceStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NotFound: Service %s not found", name)
	}
	return svc, nil
}

func (s *CloudRunService) ListServices(ctx context.Context, parent string) ([]*model.Service, error) {
	return s.serviceStore.Scan(ctx, func(k string) bool { return true })
}

func (s *CloudRunService) DeleteService(ctx context.Context, name string) error {
	return s.serviceStore.Delete(ctx, name)
}
