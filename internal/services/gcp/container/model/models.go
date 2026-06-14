package model

import "time"

type Cluster struct {
	Name            string            `json:"name"`
	Description     string            `json:"description,omitempty"`
	InitialNodeCount int               `json:"initialNodeCount"`
	Status          string            `json:"status"`
	Endpoint        string            `json:"endpoint"`
	CreateTime      time.Time         `json:"createTime"`
	Location        string            `json:"location"`
	NodePools       []*NodePool       `json:"nodePools,omitempty"`
	SelfLink        string            `json:"selfLink"`
}

type NodePool struct {
	Name             string            `json:"name"`
	Config           *NodeConfig       `json:"config"`
	InitialNodeCount int               `json:"initialNodeCount"`
	Status           string            `json:"status"`
	SelfLink         string            `json:"selfLink"`
}

type NodeConfig struct {
	MachineType string `json:"machineType"`
	DiskSizeGb  int    `json:"diskSizeGb"`
}

type ClustersList struct {
	Clusters      []*Cluster `json:"clusters"`
	NextPageToken string     `json:"nextPageToken,omitempty"`
}

type NodePoolsList struct {
	NodePools []*NodePool `json:"nodePools"`
}
