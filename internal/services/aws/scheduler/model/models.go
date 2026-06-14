package model

import (
	"time"
)

type Schedule struct {
	Name               string    `json:"Name"`
	Arn                string    `json:"Arn"`
	GroupName          string    `json:"GroupName"`
	State              string    `json:"State"` // ENABLED, DISABLED
	ScheduleExpression string    `json:"ScheduleExpression"`
	CreationDate       time.Time `json:"CreationDate"`
	LastModificationDate time.Time `json:"LastModificationDate"`
}

type ScheduleGroup struct {
	Name        string    `json:"Name"`
	Arn         string    `json:"Arn"`
	State       string    `json:"State"`
	CreationDate time.Time `json:"CreationDate"`
	LastModificationDate time.Time `json:"LastModificationDate"`
}
