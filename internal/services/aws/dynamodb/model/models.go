package model

import (
	"time"
)

type AttributeDefinition struct {
	AttributeName string `json:"AttributeName"`
	AttributeType string `json:"AttributeType"` // S, N, B
}

type KeySchemaElement struct {
	AttributeName string `json:"AttributeName"`
	KeyType       string `json:"KeyType"` // HASH, RANGE
}

type ProvisionedThroughput struct {
	ReadCapacityUnits  int64 `json:"ReadCapacityUnits"`
	WriteCapacityUnits int64 `json:"WriteCapacityUnits"`
}

type TableDefinition struct {
	TableName             string                 `json:"TableName"`
	KeySchema             []KeySchemaElement     `json:"KeySchema"`
	AttributeDefinitions  []AttributeDefinition  `json:"AttributeDefinitions"`
	TableStatus           string                 `json:"TableStatus"`
	CreationDateTime      time.Time              `json:"CreationDateTime"`
	ItemCount             int64                  `json:"ItemCount"`
	TableSizeBytes        int64                  `json:"TableSizeBytes"`
	ProvisionedThroughput *ProvisionedThroughput `json:"ProvisionedThroughput,omitempty"`
	TableArn              string                 `json:"TableArn"`
	TableID               string                 `json:"TableId,omitempty"`
	BillingMode           string                 `json:"BillingMode,omitempty"`
}
