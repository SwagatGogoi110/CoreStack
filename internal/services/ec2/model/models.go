package model

import (
	"time"
)

type Vpc struct {
	VpcID     string `json:"VpcId"`
	CidrBlock string `json:"CidrBlock"`
	State     string `json:"State"`
	IsDefault bool   `json:"IsDefault"`
	OwnerID   string `json:"OwnerId"`
}

type Subnet struct {
	SubnetID         string `json:"SubnetId"`
	VpcID            string `json:"VpcId"`
	CidrBlock        string `json:"CidrBlock"`
	State            string `json:"State"`
	AvailabilityZone string `json:"AvailabilityZone"`
}

type Instance struct {
	InstanceID   string    `json:"InstanceId"`
	ImageID      string    `json:"ImageId"`
	InstanceType string    `json:"InstanceType"`
	State        string    `json:"State"`
	LaunchTime   time.Time `json:"LaunchTime"`
	PrivateIp    string    `json:"PrivateIpAddress"`
	PublicIp     string    `json:"PublicIpAddress,omitempty"`
	VpcID        string    `json:"VpcId,omitempty"`
	SubnetID     string    `json:"SubnetId,omitempty"`
}

type Reservation struct {
	ReservationID string      `json:"ReservationId"`
	OwnerID       string      `json:"OwnerId"`
	Instances     []*Instance `json:"Instances"`
}
