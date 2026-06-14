package model

import "time"

type Pipeline struct {
	Name    string          `json:"name"`
	RoleArn string          `json:"roleArn"`
	Stages  []PipelineStage `json:"stages"`
	Version int             `json:"version"`
}

type PipelineStage struct {
	Name    string           `json:"name"`
	Actions []PipelineAction `json:"actions"`
}

type PipelineAction struct {
	Name     string         `json:"name"`
	ActionId ActionTypeId   `json:"actionTypeId"`
}

type ActionTypeId struct {
	Category string `json:"category"`
	Owner    string `json:"owner"`
	Provider string `json:"provider"`
	Version  string `json:"version"`
}

type PipelineSummary struct {
	Name             string    `json:"name"`
	Version          int       `json:"version"`
	Created          time.Time `json:"created"`
	Updated          time.Time `json:"updated"`
}
