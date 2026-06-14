package model

import "time"

type Build struct {
	Id             string            `json:"id"`
	ProjectId      string            `json:"projectId"`
	Status         string            `json:"status"`
	Source         *Source           `json:"source,omitempty"`
	Steps          []*BuildStep      `json:"steps"`
	CreateTime     time.Time         `json:"createTime"`
	StartTime      time.Time         `json:"startTime"`
	FinishTime     time.Time         `json:"finishTime"`
	Images         []string          `json:"images,omitempty"`
	LogUrl         string            `json:"logUrl"`
	Options        map[string]any    `json:"options,omitempty"`
}

type Source struct {
	StorageSource *StorageSource `json:"storageSource,omitempty"`
}

type StorageSource struct {
	Bucket string `json:"bucket"`
	Object string `json:"object"`
}

type BuildStep struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

type BuildTrigger struct {
	Id          string            `json:"id"`
	Description string            `json:"description,omitempty"`
	Filename    string            `json:"filename,omitempty"`
	CreateTime  time.Time         `json:"createTime"`
	Disabled    bool              `json:"disabled"`
}

type BuildsList struct {
	Builds        []*Build `json:"builds"`
	NextPageToken string   `json:"nextPageToken,omitempty"`
}

type TriggersList struct {
	Triggers      []*BuildTrigger `json:"triggers"`
	NextPageToken string          `json:"nextPageToken,omitempty"`
}
