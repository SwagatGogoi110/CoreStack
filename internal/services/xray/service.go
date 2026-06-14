package xray

import (
	"context"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/xray/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type XrayService struct {
	traceStore storage.Backend[string, *model.Segment]
}

func NewXrayService(factory *storage.Factory) (*XrayService, error) {
	traceStore, _ := storage.CreateAccountAware[*model.Segment](factory, "xray", "xray-traces.json", "wal")

	return &XrayService{
		traceStore: traceStore,
	}, nil
}

func (s *XrayService) PutTraceSegments(ctx context.Context, segments []string) ([]string, error) {
	var ids []string
	for _, doc := range segments {
		id := uuid.New().String()
		s.traceStore.Put(ctx, id, &model.Segment{Id: id, Document: doc})
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *XrayService) GetTraceSummaries(ctx context.Context) ([]*model.TraceSummary, error) {
	segments, _ := s.traceStore.Scan(ctx, func(k string) bool { return true })
	summaries := make([]*model.TraceSummary, 0, len(segments))
	for _, seg := range segments {
		summaries = append(summaries, &model.TraceSummary{
			Id: seg.Id,
		})
	}
	return summaries, nil
}
