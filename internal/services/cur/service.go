package cur

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/cur/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CurService struct {
	reportStore storage.Backend[string, *model.ReportDefinition]
	mu          sync.RWMutex
}

func NewCurService(factory *storage.Factory) (*CurService, error) {
	reportStore, _ := storage.CreateAccountAware[*model.ReportDefinition](factory, "cur", "cur-reports.json", "wal")

	return &CurService{
		reportStore: reportStore,
	}, nil
}

func (s *CurService) PutReportDefinition(ctx context.Context, d *model.ReportDefinition) error {
	if _, ok, _ := s.reportStore.Get(ctx, d.ReportName); ok {
		return fmt.Errorf("DuplicateReportNameException: Report already exists")
	}

	d.CreatedDate = time.Now()
	d.LastUpdatedDate = time.Now()
	d.ReportStatus = "PENDING"

	if err := s.reportStore.Put(ctx, d.ReportName, d); err != nil {
		return err
	}

	log.Printf("Created CUR report: %s", d.ReportName)
	return nil
}

func (s *CurService) DescribeReportDefinitions(ctx context.Context) ([]*model.ReportDefinition, error) {
	return s.reportStore.Scan(ctx, func(k string) bool { return true })
}
