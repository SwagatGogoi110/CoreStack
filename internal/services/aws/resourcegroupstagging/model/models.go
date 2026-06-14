package model

type ResourceTagMapping struct {
	ResourceArn string            `json:"ResourceARN"`
	Tags        map[string]string `json:"Tags"`
}

type Tag struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}
