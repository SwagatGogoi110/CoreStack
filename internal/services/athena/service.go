package athena

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/athena/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type AthenaService struct {
	queryStore storage.Backend[string, *model.QueryExecution]
}

func NewAthenaService(factory *storage.Factory) (*AthenaService, error) {
	queryStore, _ := storage.CreateAccountAware[*model.QueryExecution](factory, "athena", "athena-queries.json", "wal")

	return &AthenaService{
		queryStore: queryStore,
	}, nil
}

func (s *AthenaService) StartQueryExecution(ctx context.Context, query, workGroup string) (string, error) {
	id := uuid.New().String()
	execution := &model.QueryExecution{
		QueryExecutionID: id,
		Query:            query,
		WorkGroup:        workGroup,
		Status: model.QueryExecutionStatus{
			State:              "SUCCEEDED", // Auto-succeed in mock
			SubmissionDateTime: time.Now(),
			CompletionDateTime: time.Now(),
		},
		ResultConfiguration: model.ResultConfiguration{
			OutputLocation: fmt.Sprintf("s3://CloudStack-athena-results/%s/", id),
		},
	}

	if err := s.queryStore.Put(ctx, id, execution); err != nil {
		return "", err
	}

	log.Printf("Started Athena query: %s (%s)", query, id)
	return id, nil
}

func (s *AthenaService) GetQueryExecution(ctx context.Context, id string) (*model.QueryExecution, error) {
	execution, ok, err := s.queryStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("InvalidRequestException: Query execution not found")
	}
	return execution, nil
}
