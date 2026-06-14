package model

import "time"

type Repository struct {
	Name        string            `json:"name"`
	Format      string            `json:"format"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	CreateTime  time.Time         `json:"createTime"`
	UpdateTime  time.Time         `json:"updateTime"`
	KmsKeyName  string            `json:"kmsKeyName,omitempty"`
}

type Package struct {
	Name        string    `json:"name"`
	DisplayName string    `json:"displayName"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}

type DockerImage struct {
	Name        string    `json:"name"`
	Uri         string    `json:"uri"`
	Tags        []string  `json:"tags"`
	ImageSizeBytes int64    `json:"imageSizeBytes,string"`
	UploadTime  time.Time `json:"uploadTime"`
	MediaType   string    `json:"mediaType"`
}

type RepositoriesList struct {
	Repositories  []*Repository `json:"repositories"`
	NextPageToken string        `json:"nextPageToken,omitempty"`
}

type PackagesList struct {
	Packages      []*Package `json:"packages"`
	NextPageToken string     `json:"nextPageToken,omitempty"`
}

type DockerImagesList struct {
	DockerImages  []*DockerImage `json:"dockerImages"`
	NextPageToken string         `json:"nextPageToken,omitempty"`
}
