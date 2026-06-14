package model

import (
	"time"
)

type MetricDatum struct {
	Namespace         string       `json:"Namespace"`
	MetricName        string       `json:"MetricName"`
	Dimensions        []*Dimension `json:"Dimensions,omitempty"`
	Timestamp         time.Time    `json:"Timestamp"`
	Value             float64      `json:"Value"`
	Unit              string       `json:"Unit,omitempty"`
	SampleCount       float64      `json:"SampleCount,omitempty"`
	Sum               float64      `json:"Sum,omitempty"`
	Minimum           float64      `json:"Minimum,omitempty"`
	Maximum           float64      `json:"Maximum,omitempty"`
}

type Dimension struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

type MetricAlarm struct {
	AlarmName          string            `json:"AlarmName"`
	AlarmArn           string            `json:"AlarmArn"`
	AlarmDescription   string            `json:"AlarmDescription,omitempty"`
	MetricName         string            `json:"MetricName"`
	Namespace          string            `json:"Namespace"`
	Statistic          string            `json:"Statistic,omitempty"`
	Dimensions         []*Dimension      `json:"Dimensions,omitempty"`
	Period             int               `json:"Period"`
	EvaluationPeriods  int               `json:"EvaluationPeriods"`
	Threshold          float64           `json:"Threshold"`
	ComparisonOperator string            `json:"ComparisonOperator"`
	StateValue         string            `json:"StateValue"`
	Tags               map[string]string `json:"-"`
}
