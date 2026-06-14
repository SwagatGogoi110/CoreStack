package codepipeline

import (
	"context"
	"fmt"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/codepipeline/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CodepipelineService struct {
	pipelineStore storage.Backend[string, *model.Pipeline]
}

func NewCodepipelineService(factory *storage.Factory) (*CodepipelineService, error) {
	pipelineStore, _ := storage.CreateAccountAware[*model.Pipeline](factory, "codepipeline", "codepipeline-pipelines.json", "wal")

	return &CodepipelineService{
		pipelineStore: pipelineStore,
	}, nil
}

func (s *CodepipelineService) CreatePipeline(ctx context.Context, pipeline *model.Pipeline) (*model.Pipeline, error) {
	if _, ok, _ := s.pipelineStore.Get(ctx, pipeline.Name); ok {
		return nil, fmt.Errorf("PipelineNameInUseException: Pipeline name already in use")
	}

	pipeline.Version = 1
	if err := s.pipelineStore.Put(ctx, pipeline.Name, pipeline); err != nil {
		return nil, err
	}

	return pipeline, nil
}

func (s *CodepipelineService) GetPipeline(ctx context.Context, name string) (*model.Pipeline, error) {
	p, ok, err := s.pipelineStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("PipelineNotFoundException: Pipeline %s not found", name)
	}
	return p, nil
}

func (s *CodepipelineService) ListPipelines(ctx context.Context) ([]*model.PipelineSummary, error) {
	pipelines, _ := s.pipelineStore.Scan(ctx, func(k string) bool { return true })
	summaries := make([]*model.PipelineSummary, 0, len(pipelines))
	for _, p := range pipelines {
		summaries = append(summaries, &model.PipelineSummary{
			Name:    p.Name,
			Version: p.Version,
			Created: time.Now(), // Stub
			Updated: time.Now(), // Stub
		})
	}
	return summaries, nil
}

func (s *CodepipelineService) DeletePipeline(ctx context.Context, name string) error {
	return s.pipelineStore.Delete(ctx, name)
}
