package model

import "time"

type Trace struct {
	Name string `json:"name"`
}

type Span struct {
	Name        string    `json:"name"`
	SpanId      string    `json:"spanId"`
	ParentSpanId string   `json:"parentSpanId,omitempty"`
	DisplayName *TruncatableString `json:"displayName"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
}

type TruncatableString struct {
	Value              string `json:"value"`
	TruncatedByteCount int    `json:"truncatedByteCount"`
}

type BatchWriteSpansRequest struct {
	Spans []*Span `json:"spans"`
}
