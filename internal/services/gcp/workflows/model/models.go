package model

import "time"

type Workflow struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	State       string            `json:"state"`
	CreateTime  time.Time         `json:"createTime"`
	UpdateTime  time.Time         `json:"updateTime"`
	Labels      map[string]string `json:"labels,omitempty"`
	SourceContents string         `json:"sourceContents"`
}

type Execution struct {
	Name        string    `json:"name"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	State       string    `json:"state"`
	Argument    string    `json:"argument"`
	Result      string    `json:"result"`
	Error       *Error    `json:"error,omitempty"`
	WorkflowRevisionId string `json:"workflowRevisionId"`
}

type Error struct {
	Payload string `json:"payload"`
	Context string `json:"context"`
}

type WorkflowsList struct {
	Workflows     []*Workflow `json:"workflows"`
	NextPageToken string      `json:"nextPageToken,omitempty"`
}

type ExecutionsList struct {
	Executions    []*Execution `json:"executions"`
	NextPageToken string       `json:"nextPageToken,omitempty"`
}
