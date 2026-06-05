package model

import (
	"time"
)

type Bucket struct {
	Name             string            `json:"name"`
	CreationDate     time.Time         `json:"creationDate"`
	Region           string            `json:"region"`
	VersioningStatus string            `json:"versioningStatus,omitempty"` // "Enabled", "Suspended"
	Tags             map[string]string `json:"tags,omitempty"`
}

type S3Object struct {
	BucketName           string            `json:"bucketName"`
	Key                  string            `json:"key"`
	ContentType          string            `json:"contentType"`
	ContentEncoding      string            `json:"contentEncoding,omitempty"`
	ContentDisposition   string            `json:"contentDisposition,omitempty"`
	CacheControl         string            `json:"cacheControl,omitempty"`
	ServerSideEncryption string            `json:"serverSideEncryption,omitempty"`
	Size                 int64             `json:"size"`
	LastModified         time.Time         `json:"lastModified"`
	ETag                 string            `json:"eTag"`
	StorageClass         string            `json:"storageClass"`
	VersionID            string            `json:"versionId,omitempty"`
	IsLatest             bool              `json:"isLatest"`
	DeleteMarker         bool              `json:"deleteMarker"`
	Metadata             map[string]string `json:"metadata,omitempty"`
	Tags                 map[string]string `json:"tags,omitempty"`
}

type Part struct {
	PartNumber int    `json:"partNumber"`
	ETag       string `json:"eTag"`
	Size       int64  `json:"size"`
}

type MultipartUpload struct {
	UploadID   string            `json:"uploadId"`
	Bucket     string            `json:"bucket"`
	Key        string            `json:"key"`
	Initiated  time.Time         `json:"initiated"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	Parts      map[int]*Part     `json:"parts,omitempty"`
}
