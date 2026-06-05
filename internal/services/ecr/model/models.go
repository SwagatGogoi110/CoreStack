package model

import (
	"time"
)

type Repository struct {
	RepositoryArn      string            `json:"repositoryArn"`
	RegistryID         string            `json:"registryId"`
	RepositoryName     string            `json:"repositoryName"`
	RepositoryUri      string            `json:"repositoryUri"`
	CreatedAt          time.Time         `json:"createdAt"`
	ImageTagMutability string            `json:"imageTagMutability"` // MUTABLE | IMMUTABLE
	ScanOnPush         bool              `json:"scanOnPush"`
	Tags               map[string]string `json:"-"`
}

type ImageIdentifier struct {
	ImageTag    string `json:"imageTag,omitempty"`
	ImageDigest string `json:"imageDigest,omitempty"`
}

type ImageDetail struct {
	RegistryID     string    `json:"registryId"`
	RepositoryName string    `json:"repositoryName"`
	ImageDigest    string    `json:"imageDigest"`
	ImageTags      []string  `json:"imageTags,omitempty"`
	ImagePushedAt  time.Time `json:"imagePushedAt"`
	ImageSizeInBytes int64     `json:"imageSizeInBytes"`
}
