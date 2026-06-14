package model

import (
	"time"
)

type MskCluster struct {
	ClusterArn         string    `json:"clusterArn"`
	ClusterName        string    `json:"clusterName"`
	State              string    `json:"state"` // ACTIVE, CREATING, DELETING, FAILED
	CreationTime       time.Time `json:"creationTime"`
	CurrentVersion     string    `json:"currentVersion"`
	NumberOfBrokerNodes int       `json:"numberOfBrokerNodes"`
	ZookeeperConnectString string `json:"zookeeperConnectString"`
}
