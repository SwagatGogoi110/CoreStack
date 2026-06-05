package stepfunctions

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/stepfunctions/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type StepFunctionsService struct {
	smStore storage.Backend[string, *model.StateMachine]
	mu      sync.RWMutex
}

func NewStepFunctionsService(factory *storage.Factory) (*StepFunctionsService, error) {
	smStore, _ := storage.CreateAccountAware[*model.StateMachine](factory, "stepfunctions", "sfn-sms.json", "wal")

	return &StepFunctionsService{
		smStore: smStore,
	}, nil
}

func (s *StepFunctionsService) CreateStateMachine(ctx context.Context, name, definition, roleArn string) (*model.StateMachine, error) {
	arn := fmt.Sprintf("arn:aws:states:us-east-1:000000000000:stateMachine:%s", name)
	if _, ok, _ := s.smStore.Get(ctx, arn); ok {
		return nil, fmt.Errorf("StateMachineAlreadyExists: State machine already exists")
	}

	sm := &model.StateMachine{
		Name:            name,
		StateMachineArn: arn,
		Definition:      definition,
		RoleArn:         roleArn,
		Type:            "STANDARD",
		Status:          "ACTIVE",
		CreationDate:    time.Now(),
	}

	if err := s.smStore.Put(ctx, arn, sm); err != nil {
		return nil, err
	}

	log.Printf("Created State Machine: %s", arn)
	return sm, nil
}

func (s *StepFunctionsService) DescribeStateMachine(ctx context.Context, arn string) (*model.StateMachine, error) {
	sm, ok, err := s.smStore.Get(ctx, arn)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("StateMachineDoesNotExist: State machine does not exist")
	}
	return sm, nil
}

func (s *StepFunctionsService) ListStateMachines(ctx context.Context) ([]*model.StateMachine, error) {
	return s.smStore.Scan(ctx, func(k string) bool { return true })
}
