package model

import (
	"time"
)

type DomainStatus struct {
	DomainID      string    `json:"DomainId"`
	DomainName    string    `json:"DomainName"`
	Arn           string    `json:"ARN"`
	Created       bool      `json:"Created"`
	Deleted       bool      `json:"Deleted"`
	Endpoint      string    `json:"Endpoint,omitempty"`
	Processing    bool      `json:"Processing"`
	UpgradeProcessing bool  `json:"UpgradeProcessing"`
	EngineVersion string    `json:"EngineVersion"`
	CreatedAt     time.Time `json:"CreatedAt"`
}
