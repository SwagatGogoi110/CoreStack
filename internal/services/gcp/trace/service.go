package trace

import (
	"context"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/gcp/trace/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CloudTraceService struct {
	spanStore storage.Backend[string, *model.Span]
}

func NewCloudTraceService(factory *storage.Factory) (*CloudTraceService, error) {
	spanStore, _ := storage.CreateAccountAware[*model.Span](factory, "trace", "spans.json", "wal")

	return &CloudTraceService{
		spanStore: spanStore,
	}, nil
}

func (s *CloudTraceService) WriteSpans(ctx context.Context, spans []*model.Span) error {
	for _, span := range spans {
		s.spanStore.Put(ctx, span.Name+":"+time.Now().String(), span)
	}
	return nil
}

func (s *CloudTraceService) ListSpans(ctx context.Context) ([]*model.Span, error) {
	return s.spanStore.Scan(ctx, func(k string) bool { return true })
}
