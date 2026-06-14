package model

import "time"

type DBCluster struct {
	DBClusterIdentifier string    `json:"DBClusterIdentifier"`
	Engine              string    `json:"Engine"`
	Status              string    `json:"Status"`
	Endpoint            string    `json:"Endpoint"`
	Port                int       `json:"Port"`
	ClusterCreateTime   time.Time `json:"ClusterCreateTime"`
}
