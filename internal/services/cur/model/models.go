package model

import (
	"time"
)

type ReportDefinition struct {
	ReportName       string    `json:"ReportName"`
	TimeUnit         string    `json:"TimeUnit"`
	Format           string    `json:"Format"`
	Compression      string    `json:"Compression"`
	S3Bucket         string    `json:"S3Bucket"`
	S3Prefix         string    `json:"S3Prefix"`
	S3Region         string    `json:"S3Region"`
	CreatedDate      time.Time `json:"CreatedDate"`
	LastUpdatedDate  time.Time `json:"LastUpdatedDate"`
	ReportStatus     string    `json:"ReportStatus"`
}
