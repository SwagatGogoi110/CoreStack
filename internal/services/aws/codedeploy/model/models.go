package model

import (
	"time"
)

type Application struct {
	ApplicationID   string    `json:"applicationId"`
	ApplicationName string    `json:"applicationName"`
	CreateTime      time.Time `json:"createTime"`
	ComputePlatform string    `json:"computePlatform"`
}

type DeploymentGroup struct {
	ApplicationName      string `json:"applicationName"`
	DeploymentGroupID    string `json:"deploymentGroupId"`
	DeploymentGroupName  string `json:"deploymentGroupName"`
	DeploymentConfigName string `json:"deploymentConfigName"`
	ServiceRoleArn       string `json:"serviceRoleArn"`
}

type Deployment struct {
	DeploymentID          string    `json:"deploymentId"`
	ApplicationName       string    `json:"applicationName"`
	DeploymentGroupName   string    `json:"deploymentGroupName"`
	DeploymentConfigName  string    `json:"deploymentConfigName"`
	Status                string    `json:"status"` // Created, Queued, InProgress, Succeeded, Failed
	CreateTime            time.Time `json:"createTime"`
	StartTime             *time.Time `json:"startTime,omitempty"`
	EndTime               *time.Time `json:"endTime,omitempty"`
}
