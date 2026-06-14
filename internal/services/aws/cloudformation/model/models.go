package model

import (
	"time"
)

type Stack struct {
	StackID           string            `json:"StackId"`
	StackName         string            `json:"StackName"`
	Description       string            `json:"Description,omitempty"`
	CreationTime      time.Time         `json:"CreationTime"`
	LastUpdatedTime   *time.Time        `json:"LastUpdatedTime,omitempty"`
	StackStatus       string            `json:"StackStatus"`
	StackStatusReason string            `json:"StackStatusReason,omitempty"`
	Parameters        []*Parameter      `json:"Parameters,omitempty"`
	Outputs           []*Output         `json:"Outputs,omitempty"`
	Tags              map[string]string `json:"-"`
}

type Parameter struct {
	ParameterKey   string `json:"ParameterKey"`
	ParameterValue string `json:"ParameterValue"`
}

type Output struct {
	OutputKey   string `json:"OutputKey"`
	OutputValue string `json:"OutputValue"`
	Description string `json:"Description,omitempty"`
	ExportName  string `json:"ExportName,omitempty"`
}

type StackEvent struct {
	StackID              string    `json:"StackId"`
	EventID              string    `json:"EventId"`
	StackName            string    `json:"StackName"`
	LogicalResourceID    string    `json:"LogicalResourceId"`
	PhysicalResourceID   string    `json:"PhysicalResourceId,omitempty"`
	ResourceType         string    `json:"ResourceType"`
	Timestamp            time.Time `json:"Timestamp"`
	ResourceStatus       string    `json:"ResourceStatus"`
	ResourceStatusReason string    `json:"ResourceStatusReason,omitempty"`
}
