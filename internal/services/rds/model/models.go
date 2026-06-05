package model

import (
	"time"
)

type DbInstance struct {
	DBInstanceIdentifier string    `json:"DBInstanceIdentifier"`
	Engine               string    `json:"Engine"`
	DBInstanceStatus     string    `json:"DBInstanceStatus"` // available, creating, deleting
	DBInstanceClass      string    `json:"DBInstanceClass"`
	AllocatedStorage     int       `json:"AllocatedStorage"`
	Endpoint             *Endpoint `json:"Endpoint,omitempty"`
	InstanceCreateTime   time.Time `json:"InstanceCreateTime"`
}

type Endpoint struct {
	Address string `json:"Address"`
	Port    int    `json:"Port"`
}

type DbCluster struct {
	DBClusterIdentifier string    `json:"DBClusterIdentifier"`
	Engine              string    `json:"Engine"`
	Status              string    `json:"Status"`
	Endpoint            string    `json:"Endpoint"`
	ReaderEndpoint      string    `json:"ReaderEndpoint"`
	ClusterCreateTime   time.Time `json:"ClusterCreateTime"`
}
