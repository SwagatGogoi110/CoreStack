package apigateway

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/aws/apigateway/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

const (
	chars = "abcdefghijklmnopqrstuvwxyz0123456789"
)

type ApiGatewayService struct {
	apiStore      storage.Backend[string, *model.RestApi]
	resourceStore storage.Backend[string, *model.Resource]
}

func NewApiGatewayService(factory *storage.Factory) (*ApiGatewayService, error) {
	apiStore, _ := storage.CreateAccountAware[*model.RestApi](factory, "apigateway", "apigateway-apis.json", "wal")
	resourceStore, _ := storage.CreateAccountAware[*model.Resource](factory, "apigateway", "apigateway-resources.json", "wal")

	return &ApiGatewayService{
		apiStore:      apiStore,
		resourceStore: resourceStore,
	}, nil
}

func (s *ApiGatewayService) CreateRestApi(ctx context.Context, name, description string) (*model.RestApi, error) {
	id := s.randomID(10)
	api := &model.RestApi{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedDate: time.Now().Unix(),
	}

	if err := s.apiStore.Put(ctx, id, api); err != nil {
		return nil, err
	}

	// Create root resource
	root := &model.Resource{
		ID:   s.randomID(8),
		Path: "/",
	}
	s.resourceStore.Put(ctx, s.resourceKey(id, root.ID), root)

	log.Printf("Created REST API: %s (%s)", name, id)
	return api, nil
}

func (s *ApiGatewayService) GetRestApi(ctx context.Context, id string) (*model.RestApi, error) {
	api, ok, err := s.apiStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NotFoundException: API not found")
	}
	return api, nil
}

func (s *ApiGatewayService) ListRestApis(ctx context.Context) ([]*model.RestApi, error) {
	return s.apiStore.Scan(ctx, func(k string) bool { return true })
}

func (s *ApiGatewayService) CreateResource(ctx context.Context, apiID, parentID, pathPart string) (*model.Resource, error) {
	parent, ok, _ := s.resourceStore.Get(ctx, s.resourceKey(apiID, parentID))
	if !ok {
		return nil, fmt.Errorf("NotFoundException: Parent resource not found")
	}

	id := s.randomID(8)
	path := parent.Path
	if path == "/" {
		path = "/" + pathPart
	} else {
		path = path + "/" + pathPart
	}

	res := &model.Resource{
		ID:       id,
		ParentID: parentID,
		PathPart: pathPart,
		Path:     path,
	}

	if err := s.resourceStore.Put(ctx, s.resourceKey(apiID, id), res); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ApiGatewayService) GetResources(ctx context.Context, apiID string) ([]*model.Resource, error) {
	prefix := apiID + ":"
	return s.resourceStore.Scan(ctx, func(k string) bool {
		return len(k) > len(prefix) && k[:len(prefix)] == prefix
	})
}

func (s *ApiGatewayService) randomID(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func (s *ApiGatewayService) resourceKey(apiID, resourceID string) string {
	return fmt.Sprintf("%s:%s", apiID, resourceID)
}
