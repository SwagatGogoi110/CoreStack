package model

import "time"

type DatabaseInstance struct {
	Kind           string           `json:"kind"`
	State          string           `json:"state"`
	DatabaseVersion string          `json:"databaseVersion"`
	Settings       *Settings        `json:"settings"`
	Name           string           `json:"name"`
	Project        string           `json:"project"`
	Region         string           `json:"region"`
	IpAddresses    []*IpMapping     `json:"ipAddresses,omitempty"`
	CreateTime     time.Time        `json:"createTime"`
}

type Settings struct {
	Tier         string `json:"tier"`
	StorageAutoResize *bool `json:"storageAutoResize,omitempty"`
}

type IpMapping struct {
	Type      string `json:"type"`
	IpAddress string `json:"ipAddress"`
}

type Database struct {
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	Project  string `json:"project"`
	Instance string `json:"instance"`
	Charset  string `json:"charset,omitempty"`
	Collation string `json:"collation,omitempty"`
}

type User struct {
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
	Project  string `json:"project"`
	Instance string `json:"instance"`
	Host     string `json:"host,omitempty"`
}

type InstancesList struct {
	Kind  string              `json:"kind"`
	Items []*DatabaseInstance `json:"items"`
}

type DatabasesList struct {
	Kind  string      `json:"kind"`
	Items []*Database `json:"items"`
}

type UsersList struct {
	Kind  string  `json:"kind"`
	Items []*User `json:"items"`
}
