package bcmdataexports

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/bcmdataexports/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type BcmDataExportsService struct {
	exportStore storage.Backend[string, *model.Export]
}

func NewBcmDataExportsService(factory *storage.Factory) (*BcmDataExportsService, error) {
	exportStore, _ := storage.CreateAccountAware[*model.Export](factory, "bcmdataexports", "bcm-exports.json", "wal")

	return &BcmDataExportsService{
		exportStore: exportStore,
	}, nil
}

func (s *BcmDataExportsService) CreateExport(ctx context.Context, e *model.Export) (*model.Export, error) {
	arn := fmt.Sprintf("arn:aws:bcm-data-exports:us-east-1:000000000000:export/%s", e.Name)
	e.ExportArn = arn
	e.CreatedAt = time.Now().Unix()
	e.LastUpdatedAt = time.Now().Unix()
	e.ExportStatus = "HEALTHY"

	if err := s.exportStore.Put(ctx, arn, e); err != nil {
		return nil, err
	}

	log.Printf("Created BCM data export: %s", arn)
	return e, nil
}

func (s *BcmDataExportsService) GetExport(ctx context.Context, arn string) (*model.Export, error) {
	e, ok, err := s.exportStore.Get(ctx, arn)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("ResourceNotFoundException: Export not found")
	}
	return e, nil
}

func (s *BcmDataExportsService) ListExports(ctx context.Context) ([]*model.Export, error) {
	return s.exportStore.Scan(ctx, func(k string) bool { return true })
}
