package model

import (
	"time"
)

type CacheCluster struct {
	CacheClusterID      string    `json:"CacheClusterId"`
	CacheClusterStatus  string    `json:"CacheClusterStatus"`
	Engine              string    `json:"Engine"`
	EngineVersion       string    `json:"EngineVersion"`
	CacheClusterCreateTime time.Time `json:"CacheClusterCreateTime"`
}

type ReplicationGroup struct {
	ReplicationGroupID     string    `json:"ReplicationGroupId"`
	Description            string    `json:"Description,omitempty"`
	Status                 string    `json:"Status"` // available, creating, deleting
	AuthToken              string    `json:"-"`
	ReplicationGroupCreateTime time.Time `json:"ReplicationGroupCreateTime"`
	PrimaryEndpoint        *Endpoint `json:"PrimaryEndpoint,omitempty"`
}

type Endpoint struct {
	Address string `json:"Address"`
	Port    int    `json:"Port"`
}

type ElastiCacheUser struct {
	UserID     string    `json:"UserId"`
	UserName   string    `json:"UserName"`
	Status     string    `json:"Status"`
	CreatedTime time.Time `json:"CreatedTime"`
}
