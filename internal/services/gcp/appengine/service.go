package appengine

import (
	"context"
	"fmt"

	"github.com/hectorvent/cloudstack/internal/services/gcp/appengine/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type AppEngineService struct {
	appStore     storage.Backend[string, *model.Application]
	serviceStore storage.Backend[string, *model.Service]
	versionStore storage.Backend[string, *model.Version]
}

func NewAppEngineService(factory *storage.Factory) (*AppEngineService, error) {
	appStore, _ := storage.CreateAccountAware[*model.Application](factory, "appengine", "ae-apps.json", "wal")
	serviceStore, _ := storage.CreateAccountAware[*model.Service](factory, "appengine", "ae-services.json", "wal")
	versionStore, _ := storage.CreateAccountAware[*model.Version](factory, "appengine", "ae-versions.json", "wal")

	return &AppEngineService{
		appStore:     appStore,
		serviceStore: serviceStore,
		versionStore: versionStore,
	}, nil
}

// Applications

func (s *AppEngineService) CreateApplication(ctx context.Context, id string, app *model.Application) (*model.Application, error) {
	app.Id = id
	app.DefaultHostname = fmt.Sprintf("%s.appspot.com", id)
	if err := s.appStore.Put(ctx, id, app); err != nil {
		return nil, err
	}
	return app, nil
}

func (s *AppEngineService) GetApplication(ctx context.Context, id string) (*model.Application, error) {
	app, ok, err := s.appStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NotFound: Application %s not found", id)
	}
	return app, nil
}

// Services

func (s *AppEngineService) ListServices(ctx context.Context) ([]*model.Service, error) {
	return s.serviceStore.Scan(ctx, func(k string) bool { return true })
}

// Versions

func (s *AppEngineService) ListVersions(ctx context.Context, serviceId string) ([]*model.Version, error) {
	return s.versionStore.Scan(ctx, func(k string) bool { return true })
}
