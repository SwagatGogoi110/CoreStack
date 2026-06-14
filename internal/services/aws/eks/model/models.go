package model

import (
	"time"
)

type Cluster struct {
	Name      string            `json:"name"`
	Arn       string            `json:"arn"`
	CreatedAt time.Time         `json:"createdAt"`
	Version   string            `json:"version"`
	Endpoint  string            `json:"endpoint"`
	Status    string            `json:"status"` // CREATING, ACTIVE, DELETING, FAILED
	Tags      map[string]string `json:"tags,omitempty"`
}
