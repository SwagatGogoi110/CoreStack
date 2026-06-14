package cloudfunctions

import (
	"context"
	"fmt"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/gcp/cloudfunctions/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CloudFunctionsService struct {
	functionStore storage.Backend[string, *model.Function]
}

func NewCloudFunctionsService(factory *storage.Factory) (*CloudFunctionsService, error) {
	functionStore, _ := storage.CreateAccountAware[*model.Function](factory, "cloudfunctions", "cf-functions.json", "wal")

	return &CloudFunctionsService{
		functionStore: functionStore,
	}, nil
}

func (s *CloudFunctionsService) CreateFunction(ctx context.Context, name string, function *model.Function) (*model.Function, error) {
	if _, ok, _ := s.functionStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("AlreadyExists: Function %s already exists", name)
	}

	function.Name = name
	function.State = "ACTIVE"
	function.UpdateTime = time.Now()

	if err := s.functionStore.Put(ctx, name, function); err != nil {
		return nil, err
	}

	return function, nil
}

func (s *CloudFunctionsService) GetFunction(ctx context.Context, name string) (*model.Function, error) {
	fn, ok, err := s.functionStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NotFound: Function %s not found", name)
	}
	return fn, nil
}

func (s *CloudFunctionsService) ListFunctions(ctx context.Context, parent string) ([]*model.Function, error) {
	return s.functionStore.Scan(ctx, func(k string) bool { return true })
}

func (s *CloudFunctionsService) DeleteFunction(ctx context.Context, name string) error {
	return s.functionStore.Delete(ctx, name)
}
