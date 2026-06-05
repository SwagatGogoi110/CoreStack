package model

import (
	"time"
)

type Parameter struct {
	Name             string    `json:"Name"`
	Value            string    `json:"Value"`
	Type             string    `json:"Type"` // String, StringList, SecureString
	Version          int64     `json:"Version"`
	LastModifiedDate time.Time `json:"LastModifiedDate"`
	ARN              string    `json:"ARN"`
	DataType         string    `json:"DataType"`
}
