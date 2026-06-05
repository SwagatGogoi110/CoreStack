package model

import (
	"time"
)

type QueryExecution struct {
	QueryExecutionID    string               `json:"QueryExecutionId"`
	Query               string               `json:"Query"`
	Status              QueryExecutionStatus `json:"Status"`
	WorkGroup           string               `json:"WorkGroup,omitempty"`
	StatementType       string               `json:"StatementType,omitempty"`
	ResultConfiguration ResultConfiguration  `json:"ResultConfiguration,omitempty"`
}

type QueryExecutionStatus struct {
	State              string    `json:"State"` // QUEUED, RUNNING, SUCCEEDED, FAILED, CANCELLED
	StateChangeReason  string    `json:"StateChangeReason,omitempty"`
	SubmissionDateTime time.Time `json:"SubmissionDateTime"`
	CompletionDateTime time.Time `json:"CompletionDateTime,omitempty"`
}

type ResultConfiguration struct {
	OutputLocation string `json:"OutputLocation,omitempty"`
}

type QueryExecutionContext struct {
	Database string `json:"Database,omitempty"`
	Catalog  string `json:"Catalog,omitempty"`
}
