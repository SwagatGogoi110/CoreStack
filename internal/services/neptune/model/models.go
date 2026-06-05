package model

import (
	"time"
)

type NeptuneCluster struct {
	DBClusterIdentifier string    `json:"DBClusterIdentifier"`
	DBClusterArn        string    `json:"DBClusterArn"`
	Status              string    `json:"Status"` // available, creating, deleting
	Endpoint            string    `json:"Endpoint"`
	Port                int       `json:"Port"`
	EngineVersion       string    `json:"EngineVersion"`
	CreatedAt           time.Time `json:"CreatedAt"`
}

type NeptuneInstance struct {
	DBInstanceIdentifier string    `json:"DBInstanceIdentifier"`
	DBInstanceArn        string    `json:"DBInstanceArn"`
	DBClusterIdentifier  string    `json:"DBClusterIdentifier"`
	DBInstanceClass      string    `json:"DBInstanceClass"`
	Status              string    `json:"Status"`
	CreatedAt           time.Time `json:"CreatedAt"`
}
