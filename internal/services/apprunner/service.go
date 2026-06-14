package apprunner

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/apprunner/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type ApprunnerService struct {
	serviceStore storage.Backend[string, *model.Service]
}

func NewApprunnerService(factory *storage.Factory) (*ApprunnerService, error) {
	serviceStore, _ := storage.CreateAccountAware[*model.Service](factory, "apprunner", "apprunner-services.json", "wal")

	return &ApprunnerService{
		serviceStore: serviceStore,
	}, nil
}

func (s *ApprunnerService) CreateService(ctx context.Context, name string) (*model.Service, error) {
	if _, ok, _ := s.serviceStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("InvalidRequestException: Service %s already exists", name)
	}

	id := uuid.New().String()
	svc := &model.Service{
		ServiceArn:  fmt.Sprintf("arn:aws:apprunner:us-east-1:000000000000:service/%s/%s", name, id),
		ServiceId:   id,
		ServiceName: name,
		ServiceUrl:  fmt.Sprintf("https://%s.us-east-1.awsapprunner.com", id[:10]),
		Status:      "RUNNING",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.serviceStore.Put(ctx, name, svc)
	return svc, nil
}

func (s *ApprunnerService) DescribeService(ctx context.Context, arn string) (*model.Service, error) {
	// Simple scan for arn matching
	svcs, _ := s.serviceStore.Scan(ctx, func(k string) bool { return true })
	for _, svc := range svcs {
		if svc.ServiceArn == arn {
			return svc, nil
		}
	}
	return nil, fmt.Errorf("ResourceNotFoundException: Service with ARN %s not found", arn)
}

func (s *ApprunnerService) ListServices(ctx context.Context) ([]*model.ServiceSummary, error) {
	svcs, _ := s.serviceStore.Scan(ctx, func(k string) bool { return true })
	summaries := make([]*model.ServiceSummary, 0, len(svcs))
	for _, svc := range svcs {
		summaries = append(summaries, &model.ServiceSummary{
			ServiceArn:  svc.ServiceArn,
			ServiceId:   svc.ServiceId,
			ServiceName: svc.ServiceName,
			ServiceUrl:  svc.ServiceUrl,
			Status:      svc.Status,
			CreatedAt:   svc.CreatedAt,
			UpdatedAt:   svc.UpdatedAt,
		})
	}
	return summaries, nil
}

func (s *ApprunnerService) DeleteService(ctx context.Context, arn string) error {
	svcs, _ := s.serviceStore.Scan(ctx, func(k string) bool { return true })
	for _, svc := range svcs {
		if svc.ServiceArn == arn {
			return s.serviceStore.Delete(ctx, svc.ServiceName)
		}
	}
	return nil
}
