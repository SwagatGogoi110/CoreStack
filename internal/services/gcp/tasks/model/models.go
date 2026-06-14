package model

import "time"

type Queue struct {
	Name             string         `json:"name"`
	State            string         `json:"state"`
	RateLimits       map[string]any `json:"rateLimits,omitempty"`
	RetryConfig      map[string]any `json:"retryConfig,omitempty"`
	PurgeTime        time.Time      `json:"purgeTime,omitempty"`
}

type Task struct {
	Name             string         `json:"name"`
	HttpRequest      *HttpRequest   `json:"httpRequest,omitempty"`
	ScheduleTime     time.Time      `json:"scheduleTime"`
	CreateTime       time.Time      `json:"createTime"`
	DispatchCount    int            `json:"dispatchCount"`
	ResponseCount    int            `json:"responseCount"`
	FirstAttempt     *Attempt       `json:"firstAttempt,omitempty"`
	LastAttempt      *Attempt       `json:"lastAttempt,omitempty"`
	View             string         `json:"view"`
}

type HttpRequest struct {
	Url        string            `json:"url"`
	HttpMethod string            `json:"httpMethod"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       []byte            `json:"body,omitempty"`
}

type Attempt struct {
	ScheduleTime time.Time `json:"scheduleTime"`
	DispatchTime time.Time `json:"dispatchTime"`
	ResponseTime time.Time `json:"responseTime"`
	ResponseStatus *Status `json:"responseStatus,omitempty"`
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type QueuesList struct {
	Queues        []*Queue `json:"queues"`
	NextPageToken string   `json:"nextPageToken,omitempty"`
}

type TasksList struct {
	Tasks         []*Task `json:"tasks"`
	NextPageToken string  `json:"nextPageToken,omitempty"`
}
