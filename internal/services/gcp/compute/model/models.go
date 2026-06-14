package model

import "time"

type Instance struct {
	Kind              string            `json:"kind"`
	Id                string            `json:"id"`
	CreationTimestamp time.Time         `json:"creationTimestamp"`
	Name              string            `json:"name"`
	Description       string            `json:"description,omitempty"`
	Tags              map[string]any    `json:"tags,omitempty"`
	MachineType       string            `json:"machineType"`
	Status            string            `json:"status"`
	Zone              string            `json:"zone"`
	NetworkInterfaces []*NetworkInterface `json:"networkInterfaces"`
	Disks             []*AttachedDisk   `json:"disks"`
	SelfLink          string            `json:"selfLink"`
}

type NetworkInterface struct {
	Network   string `json:"network"`
	NetworkIP string `json:"networkIP"`
	Name      string `json:"name"`
}

type AttachedDisk struct {
	Type   string `json:"type"`
	Mode   string `json:"mode"`
	Source string `json:"source"`
	Boot   bool   `json:"boot"`
}

type Zone struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	Region    string `json:"region"`
	SelfLink  string `json:"selfLink"`
}

type Region struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	SelfLink string `json:"selfLink"`
}

type InstancesList struct {
	Kind          string      `json:"kind"`
	Id            string      `json:"id"`
	Items         []*Instance `json:"items"`
	SelfLink      string      `json:"selfLink"`
	NextPageToken string      `json:"nextPageToken,omitempty"`
}

type ZonesList struct {
	Kind  string  `json:"kind"`
	Items []*Zone `json:"items"`
}

type RegionsList struct {
	Kind  string    `json:"kind"`
	Items []*Region `json:"items"`
}
