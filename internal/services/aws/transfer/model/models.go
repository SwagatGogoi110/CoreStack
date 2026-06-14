package model

import (
	"time"
)

type Server struct {
	ServerId             string            `json:"serverId"`
	Arn                  string            `json:"arn"`
	State                string            `json:"state"` // ONLINE, OFFLINE
	Protocols            []string          `json:"protocols"`
	EndpointType         string            `json:"endpointType"`
	IdentityProviderType string            `json:"identityProviderType"`
	Tags                 map[string]string `json:"tags,omitempty"`
	CreationTime         time.Time         `json:"creationTime"`
}

type User struct {
	UserName      string            `json:"userName"`
	Arn           string            `json:"arn"`
	Role          string            `json:"role"`
	HomeDirectory string            `json:"homeDirectory"`
	Tags          map[string]string `json:"tags,omitempty"`
}
