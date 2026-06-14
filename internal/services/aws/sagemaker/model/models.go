package model

import "time"

type NotebookInstance struct {
	NotebookInstanceName string    `json:"NotebookInstanceName"`
	NotebookInstanceArn  string    `json:"NotebookInstanceArn"`
	NotebookInstanceStatus string  `json:"NotebookInstanceStatus"`
	InstanceType         string    `json:"InstanceType"`
	CreationTime         time.Time `json:"CreationTime"`
	LastModifiedTime     time.Time `json:"LastModifiedTime"`
}

type ModelSummary struct {
	ModelName    string    `json:"ModelName"`
	ModelArn     string    `json:"ModelArn"`
	CreationTime time.Time `json:"CreationTime"`
}

type EndpointSummary struct {
	EndpointName   string    `json:"EndpointName"`
	EndpointArn    string    `json:"EndpointArn"`
	EndpointStatus string    `json:"EndpointStatus"`
	CreationTime   time.Time `json:"CreationTime"`
}
