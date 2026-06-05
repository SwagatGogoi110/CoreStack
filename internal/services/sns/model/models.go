package model

import (
	"time"
)

type Topic struct {
	Name       string            `json:"name"`
	TopicArn   string            `json:"topicArn"`
	Attributes map[string]string `json:"attributes"`
	Tags       map[string]string `json:"tags"`
	CreatedAt  time.Time         `json:"createdAt"`
}

type Subscription struct {
	SubscriptionArn string            `json:"subscriptionArn"`
	TopicArn        string            `json:"topicArn"`
	Protocol        string            `json:"protocol"`
	Endpoint        string            `json:"endpoint"`
	Owner           string            `json:"owner"`
	AccountID       string            `json:"accountId"`
	Attributes      map[string]string `json:"attributes"`
}
