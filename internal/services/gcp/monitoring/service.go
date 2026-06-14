package monitoring

import (
	"context"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/gcp/monitoring/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CloudMonitoringService struct {
	tsStore storage.Backend[string, *model.TimeSeries]
}

func NewCloudMonitoringService(factory *storage.Factory) (*CloudMonitoringService, error) {
	tsStore, _ := storage.CreateAccountAware[*model.TimeSeries](factory, "monitoring", "metrics.json", "wal")

	return &CloudMonitoringService{
		tsStore: tsStore,
	}, nil
}

func (s *CloudMonitoringService) CreateTimeSeries(ctx context.Context, ts []*model.TimeSeries) error {
	for _, t := range ts {
		s.tsStore.Put(ctx, t.Metric.Type+":"+time.Now().String(), t)
	}
	return nil
}

func (s *CloudMonitoringService) ListTimeSeries(ctx context.Context) ([]*model.TimeSeries, error) {
	return s.tsStore.Scan(ctx, func(k string) bool { return true })
}
