package model

import "time"

type TimeSeries struct {
	Metric   *Metric         `json:"metric"`
	Resource *MonitoredResource `json:"resource"`
	Points   []*Point        `json:"points"`
}

type Metric struct {
	Type   string            `json:"type"`
	Labels map[string]string `json:"labels"`
}

type MonitoredResource struct {
	Type   string            `json:"type"`
	Labels map[string]string `json:"labels"`
}

type Point struct {
	Interval *Interval `json:"interval"`
	Value    *TypedValue `json:"value"`
}

type Interval struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

type TypedValue struct {
	DoubleValue *float64 `json:"doubleValue,omitempty"`
	Int64Value  *int64   `json:"int64Value,string,omitempty"`
}

type ListTimeSeriesResponse struct {
	TimeSeries    []*TimeSeries `json:"timeSeries"`
	NextPageToken string        `json:"nextPageToken,omitempty"`
}
