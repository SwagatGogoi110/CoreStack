package logging

import (
	"context"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/gcp/logging/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CloudLoggingService struct {
	entryStore storage.Backend[string, *model.LogEntry]
}

func NewCloudLoggingService(factory *storage.Factory) (*CloudLoggingService, error) {
	entryStore, _ := storage.CreateAccountAware[*model.LogEntry](factory, "logging", "log-entries.json", "wal")

	return &CloudLoggingService{
		entryStore: entryStore,
	}, nil
}

func (s *CloudLoggingService) WriteEntries(ctx context.Context, entries []*model.LogEntry) error {
	for _, e := range entries {
		if e.Timestamp.IsZero() {
			e.Timestamp = time.Now()
		}
		s.entryStore.Put(ctx, e.LogName+":"+time.Now().String(), e)
	}
	return nil
}

func (s *CloudLoggingService) ListEntries(ctx context.Context) ([]*model.LogEntry, error) {
	return s.entryStore.Scan(ctx, func(k string) bool { return true })
}
