package model

import "time"

type Secret struct {
	Name        string            `json:"name"`
	Replication *Replication      `json:"replication"`
	CreateTime  time.Time         `json:"createTime"`
	Labels      map[string]string `json:"labels,omitempty"`
	Etag        string            `json:"etag"`
}

type Replication struct {
	Automatic *Automatic `json:"automatic,omitempty"`
}

type Automatic struct {
}

type SecretVersion struct {
	Name        string    `json:"name"`
	CreateTime  time.Time `json:"createTime"`
	DestroyTime time.Time `json:"destroyTime,omitempty"`
	State       string    `json:"state"`
	Etag        string    `json:"etag"`
}

type AccessSecretVersionResponse struct {
	Name    string `json:"name"`
	Payload *Payload `json:"payload"`
}

type Payload struct {
	Data []byte `json:"data"`
}

type SecretsList struct {
	Secrets       []*Secret `json:"secrets"`
	NextPageToken string    `json:"nextPageToken,omitempty"`
}

type SecretVersionsList struct {
	Versions      []*SecretVersion `json:"versions"`
	NextPageToken string           `json:"nextPageToken,omitempty"`
}
