package model

import "time"

type Cluster struct {
	Name             string    `json:"name"`
	State            string    `json:"state"`
	CreateTime       time.Time `json:"createTime"`
	UpdateTime       time.Time `json:"updateTime"`
	CapacityConfig   *Capacity `json:"capacityConfig"`
	BootstrapAddress string    `json:"bootstrapAddress"`
}

type Capacity struct {
	VcpuCount   int64 `json:"vcpuCount"`
	MemoryBytes int64 `json:"memoryBytes"`
}

type ClustersList struct {
	Clusters      []*Cluster `json:"clusters"`
	NextPageToken string     `json:"nextPageToken,omitempty"`
}
