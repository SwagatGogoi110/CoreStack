package model

import "time"

type Job struct {
	Name             string         `json:"name"`
	Description      string         `json:"description,omitempty"`
	Schedule         string         `json:"schedule"`
	TimeZone         string         `json:"timeZone"`
	HttpTarget       *HttpTarget    `json:"httpTarget,omitempty"`
	PubsubTarget     *PubsubTarget  `json:"pubsubTarget,omitempty"`
	State            string         `json:"state"`
	Status           *Status        `json:"status"`
	LastAttemptTime  time.Time      `json:"lastAttemptTime"`
}

type HttpTarget struct {
	Uri        string            `json:"uri"`
	HttpMethod string            `json:"httpMethod"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       []byte            `json:"body,omitempty"`
}

type PubsubTarget struct {
	TopicName string            `json:"topicName"`
	Data      []byte            `json:"data"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type JobsList struct {
	Jobs          []*Job `json:"jobs"`
	NextPageToken string `json:"nextPageToken,omitempty"`
}
