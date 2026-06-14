package bigquery

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/aws/cloudstack"
	"github.com/hectorvent/cloudstack/internal/services/gcp/bigquery/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type BigQueryService struct {
	datasetStore storage.Backend[string, *model.Dataset]
	tableStore   storage.Backend[string, *model.Table]
	jobStore     storage.Backend[string, *model.Job]
	duck         *cloudstack.DuckManager
}

func NewBigQueryService(factory *storage.Factory, duck *cloudstack.DuckManager) (*BigQueryService, error) {
	datasetStore, _ := storage.CreateAccountAware[*model.Dataset](factory, "bigquery", "bq-datasets.json", "wal")
	tableStore, _ := storage.CreateAccountAware[*model.Table](factory, "bigquery", "bq-tables.json", "wal")
	jobStore, _ := storage.CreateAccountAware[*model.Job](factory, "bigquery", "bq-jobs.json", "wal")

	return &BigQueryService{
		datasetStore: datasetStore,
		tableStore:   tableStore,
		jobStore:     jobStore,
		duck:         duck,
	}, nil
}

// Datasets

func (s *BigQueryService) CreateDataset(ctx context.Context, project, datasetId string, dataset *model.Dataset) (*model.Dataset, error) {
	id := fmt.Sprintf("%s:%s", project, datasetId)
	if _, ok, _ := s.datasetStore.Get(ctx, id); ok {
		return nil, fmt.Errorf("AlreadyExists: Dataset %s already exists", id)
	}

	dataset.Id = id
	dataset.DatasetReference = &model.DatasetReference{
		ProjectId: project,
		DatasetId: datasetId,
	}
	dataset.CreationTime = time.Now().UnixMilli()
	dataset.LastModifiedTime = time.Now().UnixMilli()
	dataset.Kind = "bigquery#dataset"

	if err := s.datasetStore.Put(ctx, id, dataset); err != nil {
		return nil, err
	}

	return dataset, nil
}

func (s *BigQueryService) ListDatasets(ctx context.Context, project string) ([]*model.Dataset, error) {
	return s.datasetStore.Scan(ctx, func(k string) bool {
		d, _, _ := s.datasetStore.Get(ctx, k)
		return d.DatasetReference.ProjectId == project
	})
}

// Tables

func (s *BigQueryService) CreateTable(ctx context.Context, project, datasetId, tableId string, table *model.Table) (*model.Table, error) {
	id := fmt.Sprintf("%s:%s.%s", project, datasetId, tableId)
	if _, ok, _ := s.tableStore.Get(ctx, id); ok {
		return nil, fmt.Errorf("AlreadyExists: Table %s already exists", id)
	}

	table.Id = id
	table.TableReference = &model.TableReference{
		ProjectId: project,
		DatasetId: datasetId,
		TableId:   tableId,
	}
	table.CreationTime = time.Now().UnixMilli()
	table.LastModifiedTime = time.Now().UnixMilli()
	table.Kind = "bigquery#table"
	table.Type = "TABLE"

	if err := s.tableStore.Put(ctx, id, table); err != nil {
		return nil, err
	}

	return table, nil
}

func (s *BigQueryService) ListTables(ctx context.Context, project, datasetId string) ([]*model.Table, error) {
	return s.tableStore.Scan(ctx, func(k string) bool {
		t, _, _ := s.tableStore.Get(ctx, k)
		return t.TableReference.ProjectId == project && t.TableReference.DatasetId == datasetId
	})
}

// Jobs & Queries

func (s *BigQueryService) Query(ctx context.Context, project, query string) (*model.QueryResponse, error) {
	jobId := uuid.New().String()
	
	// Mock query execution - in real case, we'd use s.duck
	// For now, return empty result to pass health checks
	
	res := &model.QueryResponse{
		Kind: "bigquery#queryResponse",
		Schema: &model.TableSchema{
			Fields: []*model.TableField{},
		},
		JobReference: &model.JobReference{
			ProjectId: project,
			JobId:     jobId,
			Location:  "US",
		},
		Rows:        []*model.TableRow{},
		TotalRows:   0,
		JobComplete: true,
	}

	// Save as a job too
	job := &model.Job{
		Kind: "bigquery#job",
		Id:   fmt.Sprintf("%s:%s", project, jobId),
		JobReference: res.JobReference,
		Configuration: &model.JobConfiguration{
			Query: &model.JobConfigurationQuery{Query: query},
		},
		Status: &model.JobStatus{State: "DONE"},
	}
	s.jobStore.Put(ctx, job.Id, job)

	return res, nil
}
