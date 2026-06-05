package model

import (
	"time"
)

type Database struct {
	Name        string            `json:"Name"`
	Description string            `json:"Description,omitempty"`
	LocationUri string            `json:"LocationUri,omitempty"`
	Parameters  map[string]string `json:"Parameters,omitempty"`
	CreateTime  time.Time         `json:"CreateTime"`
}

type Table struct {
	Name              string            `json:"Name"`
	DatabaseName      string            `json:"DatabaseName"`
	Description       string            `json:"Description,omitempty"`
	CreateTime        time.Time         `json:"CreateTime"`
	UpdateTime        time.Time         `json:"UpdateTime"`
	StorageDescriptor *StorageDescriptor `json:"StorageDescriptor,omitempty"`
	TableType         string            `json:"TableType,omitempty"`
}

type StorageDescriptor struct {
	Columns  []*Column `json:"Columns,omitempty"`
	Location string    `json:"Location,omitempty"`
}

type Column struct {
	Name    string `json:"Name"`
	Type    string `json:"Type,omitempty"`
	Comment string `json:"Comment,omitempty"`
}
