package model

import "time"

type Cluster struct {
	ClusterIdentifier string    `json:"ClusterIdentifier"`
	NodeType          string    `json:"NodeType"`
	ClusterStatus     string    `json:"ClusterStatus"`
	MasterUsername    string    `json:"MasterUsername"`
	DBName            string    `json:"DBName"`
	Endpoint          *Endpoint `json:"Endpoint,omitempty"`
	ClusterCreateTime time.Time `json:"ClusterCreateTime"`
}

type Endpoint struct {
	Address string `json:"Address"`
	Port    int    `json:"Port"`
}
