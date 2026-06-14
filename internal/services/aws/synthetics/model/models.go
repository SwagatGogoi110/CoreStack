package model

import "time"

type Canary struct {
	Id           string        `json:"Id"`
	Name         string        `json:"Name"`
	Status       *CanaryStatus `json:"Status"`
	ArtifactS3Location string  `json:"ArtifactS3Location"`
	ExecutionRoleArn    string  `json:"ExecutionRoleArn"`
}

type CanaryStatus struct {
	State          string `json:"State"`
	StateReason    string `json:"StateReason"`
	StateReasonCode string `json:"StateReasonCode"`
}

type CanaryRun struct {
	Id       string             `json:"Id"`
	Name     string             `json:"Name"`
	Status   *CanaryRunStatus   `json:"Status"`
	Timeline *CanaryRunTimeline `json:"Timeline"`
}

type CanaryRunStatus struct {
	State       string `json:"State"`
	StateReason string `json:"StateReason"`
}

type CanaryRunTimeline struct {
	Started time.Time `json:"Started"`
	Completed time.Time `json:"Completed"`
}
