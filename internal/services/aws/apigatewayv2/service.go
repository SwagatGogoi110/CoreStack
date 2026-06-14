package apigatewayv2

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/aws/apigatewayv2/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

const (
	chars = "abcdefghijklmnopqrstuvwxyz0123456789"
)

type ApiGatewayV2Service struct {
	apiStore storage.Backend[string, *model.Api]
}

func NewApiGatewayV2Service(factory *storage.Factory) (*ApiGatewayV2Service, error) {
	apiStore, _ := storage.CreateAccountAware[*model.Api](factory, "apigatewayv2", "apigatewayv2-apis.json", "wal")

	return &ApiGatewayV2Service{
		apiStore: apiStore,
	}, nil
}

func (s *ApiGatewayV2Service) CreateApi(ctx context.Context, name, protocol, description string) (*model.Api, error) {
	id := s.randomID(10)
	api := &model.Api{
		ApiID:        id,
		Name:         name,
		ProtocolType: protocol,
		Description:  description,
		CreatedDate:  time.Now().Unix(),
		ApiEndpoint:  fmt.Sprintf("https://%s.execute-api.us-east-1.amazonaws.com", id),
	}

	if err := s.apiStore.Put(ctx, id, api); err != nil {
		return nil, err
	}

	log.Printf("Created ApiGatewayV2 API: %s (%s)", name, id)
	return api, nil
}

func (s *ApiGatewayV2Service) GetApi(ctx context.Context, id string) (*model.Api, error) {
	api, ok, err := s.apiStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NotFoundException: API not found")
	}
	return api, nil
}

func (s *ApiGatewayV2Service) GetApis(ctx context.Context) ([]*model.Api, error) {
	return s.apiStore.Scan(ctx, func(k string) bool { return true })
}

func (s *ApiGatewayV2Service) randomID(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
