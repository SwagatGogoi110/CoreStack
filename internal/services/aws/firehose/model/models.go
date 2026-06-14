package model

import (
	"time"
)

type DeliveryStreamDescription struct {
	DeliveryStreamName   string                `json:"DeliveryStreamName"`
	DeliveryStreamArn    string                `json:"DeliveryStreamARN"`
	DeliveryStreamStatus string                `json:"DeliveryStreamStatus"`
	CreateTimestamp      time.Time             `json:"CreateTimestamp"`
	Destinations         []*DestinationDetails `json:"Destinations"`
}

type DestinationDetails struct {
	S3DestinationDescription *S3DestinationDescription `json:"S3DestinationDescription,omitempty"`
}

type S3DestinationDescription struct {
	BucketArn string `json:"BucketARN"`
	Prefix    string `json:"Prefix,omitempty"`
}

type Record struct {
	Data []byte `json:"Data"`
}
