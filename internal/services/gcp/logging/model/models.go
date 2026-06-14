package model

import "time"

type LogEntry struct {
	LogName          string            `json:"logName"`
	Resource         *MonitoredResource `json:"resource"`
	Timestamp        time.Time         `json:"timestamp"`
	Severity         string            `json:"severity"`
	TextPayload      string            `json:"textPayload,omitempty"`
	JsonPayload      map[string]any    `json:"jsonPayload,omitempty"`
	Labels           map[string]string `json:"labels,omitempty"`
}

type MonitoredResource struct {
	Type   string            `json:"type"`
	Labels map[string]string `json:"labels"`
}

type WriteLogEntriesRequest struct {
	LogName  string      `json:"logName,omitempty"`
	Entries  []*LogEntry `json:"entries"`
	Labels   map[string]string `json:"labels,omitempty"`
}

type ListLogEntriesResponse struct {
	Entries       []*LogEntry `json:"entries"`
	NextPageToken string      `json:"nextPageToken,omitempty"`
}
