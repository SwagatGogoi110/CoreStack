package model

import (
	"time"
)

type Project struct {
	Name        string    `json:"name"`
	Arn         string    `json:"arn"`
	Description string    `json:"description,omitempty"`
	Created     time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
}

type Build struct {
	ID           string    `json:"id"`
	Arn          string    `json:"arn"`
	BuildNumber  int64     `json:"buildNumber"`
	BuildStatus  string    `json:"buildStatus"` // IN_PROGRESS, SUCCEEDED, FAILED
	ProjectName  string    `json:"projectName"`
	StartTime    time.Time `json:"startTime"`
	EndTime      *time.Time `json:"endTime,omitempty"`
}
