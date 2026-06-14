package model

import (
	"time"
)

type LoadBalancer struct {
	LoadBalancerArn  string    `json:"LoadBalancerArn"`
	DNSName          string    `json:"DNSName"`
	CreatedTime      time.Time `json:"CreatedTime"`
	LoadBalancerName string    `json:"LoadBalancerName"`
	Scheme           string    `json:"Scheme"`
	VpcID            string    `json:"VpcId"`
	State            string    `json:"State"`
	Type             string    `json:"Type"`
}

type TargetGroup struct {
	TargetGroupArn  string `json:"TargetGroupArn"`
	TargetGroupName string `json:"TargetGroupName"`
	Protocol        string `json:"Protocol"`
	Port            int    `json:"Port"`
	VpcID           string `json:"VpcId"`
	TargetType      string `json:"TargetType"`
}

type Listener struct {
	ListenerArn     string `json:"ListenerArn"`
	LoadBalancerArn string `json:"LoadBalancerArn"`
	Port            int    `json:"Port"`
	Protocol        string `json:"Protocol"`
}
