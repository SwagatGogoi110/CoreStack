package model

import "time"

type Topic struct {
	Name                     string            `json:"name"`
	Labels                   map[string]string `json:"labels,omitempty"`
	MessageRetentionDuration string            `json:"messageRetentionDuration,omitempty"`
}

type Subscription struct {
	Name                     string         `json:"name"`
	Topic                    string         `json:"topic"`
	PushConfig               *PushConfig    `json:"pushConfig,omitempty"`
	AckDeadlineSeconds       int            `json:"ackDeadlineSeconds"`
	RetainAckedMessages      bool           `json:"retainAckedMessages"`
	MessageRetentionDuration string         `json:"messageRetentionDuration"`
	ExpirationPolicy         map[string]any `json:"expirationPolicy,omitempty"`
	Filter                   string         `json:"filter,omitempty"`
	DeadLetterPolicy         map[string]any `json:"deadLetterPolicy,omitempty"`
	RetryPolicy              map[string]any `json:"retryPolicy,omitempty"`
	Detached                 bool           `json:"detached"`
}

type PushConfig struct {
	PushEndpoint string            `json:"pushEndpoint"`
	Attributes   map[string]string `json:"attributes,omitempty"`
}

type Message struct {
	Data        []byte            `json:"data"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	MessageId   string            `json:"messageId"`
	PublishTime time.Time         `json:"publishTime"`
	OrderingKey string            `json:"orderingKey,omitempty"`
}

type ReceivedMessage struct {
	AckId   string   `json:"ackId"`
	Message *Message `json:"message"`
}

type TopicsList struct {
	Topics        []*Topic `json:"topics"`
	NextPageToken string   `json:"nextPageToken,omitempty"`
}

type SubscriptionsList struct {
	Subscriptions []*Subscription `json:"subscriptions"`
	NextPageToken string          `json:"nextPageToken,omitempty"`
}
