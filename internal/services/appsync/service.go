package appsync

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/hectorvent/cloudstack/internal/services/appsync/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

const (
	chars = "abcdefghijklmnopqrstuvwxyz0123456789"
)

type AppSyncService struct {
	apiStore storage.Backend[string, *model.GraphqlApi]
}

func NewAppSyncService(factory *storage.Factory) (*AppSyncService, error) {
	apiStore, _ := storage.CreateAccountAware[*model.GraphqlApi](factory, "appsync", "appsync-apis.json", "wal")

	return &AppSyncService{
		apiStore: apiStore,
	}, nil
}

func (s *AppSyncService) CreateGraphqlApi(ctx context.Context, name, authType string) (*model.GraphqlApi, error) {
	id := s.randomID(26)
	api := &model.GraphqlApi{
		ApiID:             id,
		Name:              name,
		Arn:               fmt.Sprintf("arn:aws:appsync:us-east-1:000000000000:apis/%s", id),
		AuthenticationType: authType,
		Uris: map[string]string{
			"GRAPHQL":  fmt.Sprintf("http://localhost:8080/graphql/%s", id),
			"REALTIME": fmt.Sprintf("ws://localhost:8080/graphql/%s", id),
		},
	}

	if err := s.apiStore.Put(ctx, id, api); err != nil {
		return nil, err
	}

	log.Printf("Created AppSync GraphQL API: %s (%s)", name, id)
	return api, nil
}

func (s *AppSyncService) GetGraphqlApi(ctx context.Context, id string) (*model.GraphqlApi, error) {
	api, ok, err := s.apiStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NotFoundException: GraphQL API not found")
	}
	return api, nil
}

func (s *AppSyncService) ListGraphqlApis(ctx context.Context) ([]*model.GraphqlApi, error) {
	return s.apiStore.Scan(ctx, func(k string) bool { return true })
}

func (s *AppSyncService) randomID(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
