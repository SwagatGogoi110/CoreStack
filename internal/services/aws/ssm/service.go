package ssm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/aws/ssm/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type SsmService struct {
	parameterStore storage.Backend[string, *model.Parameter]
}

func NewSsmService(factory *storage.Factory) (*SsmService, error) {
	parameterStore, _ := storage.CreateAccountAware[*model.Parameter](factory, "ssm", "ssm-parameters.json", "wal")

	return &SsmService{
		parameterStore: parameterStore,
	}, nil
}

func (s *SsmService) PutParameter(ctx context.Context, name, value, paramType string, overwrite bool) (int64, error) {
	existing, ok, _ := s.parameterStore.Get(ctx, name)
	if ok && !overwrite {
		return 0, fmt.Errorf("ParameterAlreadyExists: The parameter already exists")
	}

	version := int64(1)
	if ok {
		version = existing.Version + 1
	}

	param := &model.Parameter{
		Name:             name,
		Value:            value,
		Type:             paramType,
		Version:          version,
		LastModifiedDate: time.Now(),
		ARN:              fmt.Sprintf("arn:aws:ssm:us-east-1:000000000000:parameter%s", name),
		DataType:         "text",
	}

	if err := s.parameterStore.Put(ctx, name, param); err != nil {
		return 0, err
	}

	log.Printf("PutParameter: %s (version %d)", name, version)
	return version, nil
}

func (s *SsmService) GetParameter(ctx context.Context, name string) (*model.Parameter, error) {
	param, ok, err := s.parameterStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("ParameterNotFound: Parameter not found")
	}
	return param, nil
}

func (s *SsmService) GetParameters(ctx context.Context, names []string) ([]*model.Parameter, error) {
	var result []*model.Parameter
	for _, n := range names {
		if p, ok, _ := s.parameterStore.Get(ctx, n); ok {
			result = append(result, p)
		}
	}
	return result, nil
}
