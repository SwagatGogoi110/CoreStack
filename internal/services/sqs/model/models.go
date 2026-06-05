package model

import (
	"time"
)

type Queue struct {
	QueueName             string            `json:"queueName"`
	QueueURL              string            `json:"queueUrl"`
	AccountID             string            `json:"accountId"`
	Attributes            map[string]string `json:"attributes"`
	Tags                  map[string]string `json:"tags"`
	CreatedTimestamp      time.Time         `json:"createdTimestamp"`
	LastModifiedTimestamp time.Time         `json:"lastModifiedTimestamp"`
}

type Message struct {
	MessageID              string                           `json:"messageId"`
	Body                   string                           `json:"body"`
	MessageAttributes      map[string]*MessageAttributeValue `json:"messageAttributes,omitempty"`
	SentTimestamp          time.Time                        `json:"sentTimestamp"`
	FirstReceiveTimestamp  *time.Time                       `json:"firstReceiveTimestamp,omitempty"`
	ReceiveCount           int                              `json:"receiveCount"`
	MD5OfBody              string                           `json:"md5OfBody"`
	MD5OfMessageAttributes string                           `json:"md5OfMessageAttributes,omitempty"`
	ReceiptHandle          string                           `json:"-"`
	VisibleAt              time.Time                        `json:"-"`
}

type MessageAttributeValue struct {
	StringValue string `json:"stringValue,omitempty"`
	BinaryValue []byte `json:"binaryValue,omitempty"`
	DataType    string `json:"dataType"`
}
