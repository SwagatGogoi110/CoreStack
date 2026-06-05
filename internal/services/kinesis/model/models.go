package model

import (
	"time"
)

type KinesisStream struct {
	StreamName              string            `json:"StreamName"`
	StreamArn               string            `json:"StreamARN"`
	StreamStatus            string            `json:"StreamStatus"`
	Shards                  []*KinesisShard   `json:"Shards"`
	RetentionPeriodHours    int               `json:"RetentionPeriodHours"`
	StreamCreationTimestamp time.Time         `json:"StreamCreationTimestamp"`
	Tags                    map[string]string `json:"Tags,omitempty"`
}

type KinesisShard struct {
	ShardID string `json:"ShardId"`
}

type KinesisRecord struct {
	Data                        []byte    `json:"Data"`
	PartitionKey                string    `json:"PartitionKey"`
	SequenceNumber              string    `json:"SequenceNumber"`
	ApproximateArrivalTimestamp time.Time `json:"ApproximateArrivalTimestamp"`
}
