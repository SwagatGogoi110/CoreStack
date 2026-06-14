package model

import "time"

type Service struct {
	ServiceArn    string         `json:"ServiceArn"`
	ServiceId     string         `json:"ServiceId"`
	ServiceName   string         `json:"ServiceName"`
	ServiceUrl    string         `json:"ServiceUrl"`
	Status        string         `json:"Status"`
	CreatedAt     time.Time      `json:"CreatedAt"`
	UpdatedAt     time.Time      `json:"UpdatedAt"`
}

type ServiceSummary struct {
	ServiceArn  string    `json:"ServiceArn"`
	ServiceId   string    `json:"ServiceId"`
	ServiceName string    `json:"ServiceName"`
	ServiceUrl  string    `json:"ServiceUrl"`
	Status      string    `json:"Status"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}
