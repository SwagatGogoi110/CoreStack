package model

import (
	"time"
)

type EventBus struct {
	Name        string            `json:"name"`
	Arn         string            `json:"arn"`
	Description string            `json:"description,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
	CreatedTime time.Time         `json:"createdTime"`
}

type Rule struct {
	Name               string            `json:"name"`
	Arn                string            `json:"arn"`
	EventBusName       string            `json:"eventBusName"`
	EventPattern       string            `json:"eventPattern,omitempty"`
	ScheduleExpression string            `json:"scheduleExpression,omitempty"`
	State              string            `json:"state"` // ENABLED, DISABLED
	Description        string            `json:"description,omitempty"`
	Tags               map[string]string `json:"tags,omitempty"`
	CreatedAt          time.Time         `json:"createdAt"`
}

type Target struct {
	ID   string `json:"id"`
	Arn  string `json:"arn"`
	Input string `json:"input,omitempty"`
}
