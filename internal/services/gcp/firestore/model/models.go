package model

import "time"

type Document struct {
	Name       string           `json:"name"`
	Fields     map[string]Value `json:"fields,omitempty"`
	CreateTime time.Time        `json:"createTime"`
	UpdateTime time.Time        `json:"updateTime"`
}

type Value struct {
	NullValue     any            `json:"nullValue,omitempty"`
	BooleanValue  *bool          `json:"booleanValue,omitempty"`
	IntegerValue  *string        `json:"integerValue,omitempty"`
	DoubleValue   *float64       `json:"doubleValue,omitempty"`
	TimestampValue *time.Time    `json:"timestampValue,omitempty"`
	StringValue   *string        `json:"stringValue,omitempty"`
	BytesValue    []byte         `json:"bytesValue,omitempty"`
	ReferenceValue *string       `json:"referenceValue,omitempty"`
	GeoPointValue *map[string]any `json:"geoPointValue,omitempty"`
	ArrayValue    *ArrayValue    `json:"arrayValue,omitempty"`
	MapValue      *MapValue      `json:"mapValue,omitempty"`
}

type ArrayValue struct {
	Values []Value `json:"values"`
}

type MapValue struct {
	Fields map[string]Value `json:"fields"`
}

type DocumentsList struct {
	Documents     []*Document `json:"documents"`
	NextPageToken string      `json:"nextPageToken,omitempty"`
}
