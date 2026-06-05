package model

import (
	"time"
)

type Secret struct {
	Name             string                    `json:"Name"`
	Arn              string                    `json:"ARN"`
	Description      string                    `json:"Description,omitempty"`
	KmsKeyId         string                    `json:"KmsKeyId,omitempty"`
	RotationEnabled  bool                      `json:"RotationEnabled"`
	CreatedDate      time.Time                 `json:"CreatedDate"`
	LastChangedDate  time.Time                 `json:"LastChangedDate"`
	Versions         map[string]*SecretVersion `json:"-"`
	CurrentVersionId string                    `json:"-"`
}

type SecretVersion struct {
	VersionId     string    `json:"VersionId"`
	SecretString  string    `json:"SecretString,omitempty"`
	SecretBinary  []byte    `json:"SecretBinary,omitempty"`
	VersionStages []string  `json:"VersionStages"`
	CreatedDate   time.Time `json:"CreatedDate"`
}
