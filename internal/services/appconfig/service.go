package appconfig

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/hectorvent/cloudstack/internal/services/appconfig/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

const (
	chars = "abcdefghijklmnopqrstuvwxyz0123456789"
)

type AppConfigService struct {
	appStore storage.Backend[string, *model.Application]
}

func NewAppConfigService(factory *storage.Factory) (*AppConfigService, error) {
	appStore, _ := storage.CreateAccountAware[*model.Application](factory, "appconfig", "appconfig-apps.json", "wal")

	return &AppConfigService{
		appStore: appStore,
	}, nil
}

func (s *AppConfigService) CreateApplication(ctx context.Context, name, description string) (*model.Application, error) {
	id := s.randomID(7)
	app := &model.Application{
		ID:          id,
		Name:        name,
		Description: description,
	}

	if err := s.appStore.Put(ctx, id, app); err != nil {
		return nil, err
	}

	log.Printf("Created AppConfig application: %s (%s)", name, id)
	return app, nil
}

func (s *AppConfigService) GetApplication(ctx context.Context, id string) (*model.Application, error) {
	app, ok, err := s.appStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("ResourceNotFoundException: Application not found")
	}
	return app, nil
}

func (s *AppConfigService) ListApplications(ctx context.Context) ([]*model.Application, error) {
	return s.appStore.Scan(ctx, func(k string) bool { return true })
}

func (s *AppConfigService) randomID(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
