package model

import "time"

type Cluster struct {
	Id              string         `json:"Id"`
	Name            string         `json:"Name"`
	Status          *ClusterStatus `json:"Status"`
	ReleaseLabel    string         `json:"ReleaseLabel"`
	NormalizedInstanceHours int    `json:"NormalizedInstanceHours"`
}

type ClusterStatus struct {
	State          string                `json:"State"`
	StateChangeReason *StateChangeReason `json:"StateChangeReason,omitempty"`
	Timeline       *Timeline             `json:"Timeline"`
}

type StateChangeReason struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

type Timeline struct {
	CreationDateTime time.Time `json:"CreationDateTime"`
	ReadyDateTime    time.Time `json:"ReadyDateTime"`
}
