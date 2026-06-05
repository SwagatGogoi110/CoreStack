package model

import (
	"time"
)

type LambdaFunction struct {
	FunctionName    string            `json:"functionName"`
	FunctionArn     string            `json:"functionArn"`
	Runtime         string            `json:"runtime"`
	Role            string            `json:"role"`
	Handler         string            `json:"handler"`
	Description     string            `json:"description,omitempty"`
	Timeout         int               `json:"timeout"`
	MemorySize      int               `json:"memorySize"`
	CodeSizeBytes   int64             `json:"codeSizeBytes"`
	Environment     map[string]string `json:"environment,omitempty"`
	Tags            map[string]string `json:"tags,omitempty"`
	LastModified    time.Time         `json:"lastModified"`
	RevisionID      string            `json:"revisionId"`
	CodeLocalPath   string            `json:"codeLocalPath,omitempty"`
}

type InvokeResult struct {
	StatusCode int    `json:"statusCode"`
	Payload    []byte `json:"payload"`
	FunctionError string `json:"functionError,omitempty"`
}
