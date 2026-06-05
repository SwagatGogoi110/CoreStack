package model

import (
	"time"
)

type EcsCluster struct {
	ClusterName       string            `json:"clusterName"`
	ClusterArn        string            `json:"clusterArn"`
	Status            string            `json:"status"` // ACTIVE, INACTIVE
	RegisteredTasks   int               `json:"registeredContainerInstancesCount"`
	RunningTasksCount int               `json:"runningTasksCount"`
	Tags              map[string]string `json:"tags,omitempty"`
}

type TaskDefinition struct {
	Family            string `json:"family"`
	Revision          int    `json:"revision"`
	TaskDefinitionArn string `json:"taskDefinitionArn"`
	Status            string `json:"status"` // ACTIVE, INACTIVE
}

type EcsServiceModel struct {
	ServiceName string `json:"serviceName"`
	ServiceArn  string `json:"serviceArn"`
	ClusterArn  string `json:"clusterArn"`
	Status      string `json:"status"` // ACTIVE, DRAINING, INACTIVE
}

type EcsTask struct {
	TaskArn           string    `json:"taskArn"`
	ClusterArn        string    `json:"clusterArn"`
	TaskDefinitionArn string    `json:"taskDefinitionArn"`
	LastStatus        string    `json:"lastStatus"`
	DesiredStatus     string    `json:"desiredStatus"`
	CreatedAt         time.Time `json:"createdAt"`
}
