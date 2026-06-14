package model

import "time"

type Repository struct {
	AccountId      string    `json:"accountId"`
	RepositoryId   string    `json:"repositoryId"`
	RepositoryName string    `json:"repositoryName"`
	RepositoryDescription string `json:"repositoryDescription"`
	LastModifiedDate time.Time `json:"lastModifiedDate"`
	CreationDate     time.Time `json:"creationDate"`
	CloneUrlHttp     string    `json:"cloneUrlHttp"`
	CloneUrlSsh      string    `json:"cloneUrlSsh"`
	Arn              string    `json:"Arn"`
}
