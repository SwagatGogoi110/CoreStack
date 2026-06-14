package model

import (
	"time"
)

type AutoScalingGroup struct {
	AutoScalingGroupName    string            `json:"AutoScalingGroupName"`
	AutoScalingGroupArn     string            `json:"AutoScalingGroupARN"`
	LaunchConfigurationName string            `json:"LaunchConfigurationName,omitempty"`
	MinSize                 int               `json:"MinSize"`
	MaxSize                 int               `json:"MaxSize"`
	DesiredCapacity         int               `json:"DesiredCapacity"`
	DefaultCooldown         int               `json:"DefaultCooldown"`
	AvailabilityZones       []string          `json:"AvailabilityZones"`
	CreatedTime             time.Time         `json:"CreatedTime"`
	Tags                    map[string]string `json:"-"`
}

type LaunchConfiguration struct {
	LaunchConfigurationName string    `json:"LaunchConfigurationName"`
	LaunchConfigurationArn  string    `json:"LaunchConfigurationARN"`
	ImageID                 string    `json:"ImageId"`
	InstanceType            string    `json:"InstanceType"`
	KeyName                 string    `json:"KeyName,omitempty"`
	CreatedTime             time.Time `json:"CreatedTime"`
}
