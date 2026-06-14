package model

import "time"

type Function struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	BuildConfig *BuildConfig      `json:"buildConfig,omitempty"`
	ServiceConfig *ServiceConfig   `json:"serviceConfig,omitempty"`
	State       string            `json:"state"`
	UpdateTime  time.Time         `json:"updateTime"`
	Labels      map[string]string `json:"labels,omitempty"`
}

type BuildConfig struct {
	Runtime     string `json:"runtime"`
	EntryPoint  string `json:"entryPoint"`
	Source      *Source `json:"source,omitempty"`
}

type Source struct {
	StorageSource *StorageSource `json:"storageSource,omitempty"`
}

type StorageSource struct {
	Bucket string `json:"bucket"`
	Object string `json:"object"`
}

type ServiceConfig struct {
	Service       string `json:"service"`
	TimeoutSeconds int    `json:"timeoutSeconds"`
	AvailableMemory string `json:"availableMemory"`
}

type FunctionsList struct {
	Functions     []*Function `json:"functions"`
	NextPageToken string      `json:"nextPageToken,omitempty"`
}
