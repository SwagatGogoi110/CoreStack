package model

import "time"

type Instance struct {
	Name            string            `json:"name"`
	Config          string            `json:"config"`
	DisplayName     string            `json:"displayName"`
	NodeCount       int               `json:"nodeCount,omitempty"`
	ProcessingUnits int               `json:"processingUnits,omitempty"`
	State           string            `json:"state"`
	Labels          map[string]string `json:"labels,omitempty"`
	CreateTime      time.Time         `json:"createTime"`
	UpdateTime      time.Time         `json:"updateTime"`
}

type Database struct {
	Name       string    `json:"name"`
	State      string    `json:"state"`
	CreateTime time.Time `json:"createTime"`
}

type Session struct {
	Name                   string    `json:"name"`
	ApproximateLastUseTime time.Time `json:"approximateLastUseTime"`
	CreateTime             time.Time `json:"createTime"`
}

type InstancesList struct {
	Instances     []*Instance `json:"instances"`
	NextPageToken string      `json:"nextPageToken,omitempty"`
}

type DatabasesList struct {
	Databases     []*Database `json:"databases"`
	NextPageToken string      `json:"nextPageToken,omitempty"`
}

type ResultSet struct {
	Metadata *ResultSetMetadata `json:"metadata"`
	Rows     [][]any           `json:"rows"`
	Stats    map[string]any    `json:"stats,omitempty"`
}

type ResultSetMetadata struct {
	RowType *StructType `json:"rowType"`
}

type StructType struct {
	Fields []*Field `json:"fields"`
}

type Field struct {
	Name string `json:"name"`
	Type *Type  `json:"type"`
}

type Type struct {
	Code string `json:"code"`
}
