package cloudformation

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/cloudformation/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CloudFormationService struct {
	stackStore storage.Backend[string, *model.Stack]
}

func NewCloudFormationService(factory *storage.Factory) (*CloudFormationService, error) {
	stackStore, _ := storage.CreateAccountAware[*model.Stack](factory, "cloudformation", "cf-stacks.json", "wal")

	return &CloudFormationService{
		stackStore: stackStore,
	}, nil
}

func (s *CloudFormationService) CreateStack(ctx context.Context, name, template string) (*model.Stack, error) {
	if _, ok, _ := s.stackStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("AlreadyExistsException: Stack already exists")
	}

	id := fmt.Sprintf("arn:aws:cloudformation:us-east-1:000000000000:stack/%s/%s", name, uuid.New().String())
	stack := &model.Stack{
		StackID:      id,
		StackName:    name,
		StackStatus:  "CREATE_COMPLETE", // Auto-complete for now
		CreationTime: time.Now(),
	}

	if err := s.stackStore.Put(ctx, name, stack); err != nil {
		return nil, err
	}

	log.Printf("Created CloudFormation stack: %s", name)
	return stack, nil
}

func (s *CloudFormationService) DescribeStacks(ctx context.Context, name string) ([]*model.Stack, error) {
	if name != "" {
		stack, ok, err := s.stackStore.Get(ctx, name)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("ValidationError: Stack does not exist")
		}
		return []*model.Stack{stack}, nil
	}
	return s.stackStore.Scan(ctx, func(k string) bool { return true })
}

func (s *CloudFormationService) DeleteStack(ctx context.Context, name string) error {
	return s.stackStore.Delete(ctx, name)
}
