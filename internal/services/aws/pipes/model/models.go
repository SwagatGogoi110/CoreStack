package model

import (
	"time"
)

type Pipe struct {
	Name             string    `json:"Name"`
	Arn              string    `json:"Arn"`
	Source           string    `json:"Source"`
	Target           string    `json:"Target"`
	State            string    `json:"CurrentState"` // RUNNING, STOPPED
	DesiredState     string    `json:"DesiredState"` // RUNNING, STOPPED
	CreationTime     time.Time `json:"CreationTime"`
	LastModifiedTime time.Time `json:"LastModifiedTime"`
}
