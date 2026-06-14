package model

import (
	"time"
)

type Distribution struct {
	ID                 string              `json:"Id"`
	Arn                string              `json:"ARN"`
	Status             string              `json:"Status"`
	LastModifiedTime   time.Time           `json:"LastModifiedTime"`
	DomainName         string              `json:"DomainName"`
	DistributionConfig *DistributionConfig `json:"DistributionConfig"`
}

type DistributionConfig struct {
	CallerReference string  `json:"CallerReference"`
	Aliases         *Aliases `json:"Aliases,omitempty"`
	DefaultRootObject string `json:"DefaultRootObject,omitempty"`
	Origins         *Origins `json:"Origins"`
	Enabled         bool    `json:"Enabled"`
}

type Aliases struct {
	Quantity int      `json:"Quantity"`
	Items    []string `json:"Items,omitempty"`
}

type Origins struct {
	Quantity int       `json:"Quantity"`
	Items    []*Origin `json:"Items"`
}

type Origin struct {
	ID         string `json:"Id"`
	DomainName string `json:"DomainName"`
}
