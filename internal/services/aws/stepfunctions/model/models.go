package model

import (
	"time"
)

type StateMachine struct {
	Name            string    `json:"name"`
	StateMachineArn string    `json:"stateMachineArn"`
	Definition      string    `json:"definition"`
	RoleArn         string    `json:"roleArn"`
	Type            string    `json:"type"` // STANDARD, EXPRESS
	Status          string    `json:"status"`
	CreationDate    time.Time `json:"creationDate"`
}

type Execution struct {
	ExecutionArn    string    `json:"executionArn"`
	StateMachineArn string    `json:"stateMachineArn"`
	Name            string    `json:"name"`
	Status          string    `json:"status"` // RUNNING, SUCCEEDED, FAILED, TIMED_OUT, ABORTED
	StartDate       time.Time `json:"startDate"`
	StopDate        *time.Time `json:"stopDate,omitempty"`
	Input           string    `json:"input,omitempty"`
	Output          string    `json:"output,omitempty"`
}
