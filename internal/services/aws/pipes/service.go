package pipes

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/aws/pipes/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type PipesService struct {
	pipeStore storage.Backend[string, *model.Pipe]
	mu        sync.RWMutex
}

func NewPipesService(factory *storage.Factory) (*PipesService, error) {
	pipeStore, _ := storage.CreateAccountAware[*model.Pipe](factory, "pipes", "pipes.json", "wal")

	return &PipesService{
		pipeStore: pipeStore,
	}, nil
}

func (s *PipesService) CreatePipe(ctx context.Context, name, source, target string) (*model.Pipe, error) {
	if _, ok, _ := s.pipeStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("ConflictException: Pipe %s already exists", name)
	}

	arn := fmt.Sprintf("arn:aws:pipes:us-east-1:000000000000:pipe/%s", name)
	pipe := &model.Pipe{
		Name:             name,
		Arn:              arn,
		Source:           source,
		Target:           target,
		State:            "RUNNING",
		DesiredState:     "RUNNING",
		CreationTime:     time.Now(),
		LastModifiedTime: time.Now(),
	}

	if err := s.pipeStore.Put(ctx, name, pipe); err != nil {
		return nil, err
	}

	log.Printf("Created Pipe: %s", name)
	return pipe, nil
}

func (s *PipesService) DescribePipe(ctx context.Context, name string) (*model.Pipe, error) {
	pipe, ok, err := s.pipeStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NotFoundException: Pipe not found")
	}
	return pipe, nil
}

func (s *PipesService) ListPipes(ctx context.Context) ([]*model.Pipe, error) {
	return s.pipeStore.Scan(ctx, func(k string) bool { return true })
}
