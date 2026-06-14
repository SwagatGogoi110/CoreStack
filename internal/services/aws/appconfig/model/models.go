package model

type Application struct {
	ID          string            `json:"Id"`
	Name        string            `json:"Name"`
	Description string            `json:"Description,omitempty"`
	Tags        map[string]string `json:"Tags,omitempty"`
}

type Environment struct {
	ID            string `json:"Id"`
	ApplicationID string `json:"ApplicationId"`
	Name          string `json:"Name"`
	Description   string `json:"Description,omitempty"`
	State         string `json:"State"`
}

type ConfigurationProfile struct {
	ID            string `json:"Id"`
	ApplicationID string `json:"ApplicationId"`
	Name          string `json:"Name"`
	Description   string `json:"Description,omitempty"`
	LocationUri   string `json:"LocationUri"`
	Type          string `json:"Type"`
}

type HostedConfigurationVersion struct {
	ApplicationID          string `json:"ApplicationId"`
	ConfigurationProfileID string `json:"ConfigurationProfileId"`
	VersionNumber          int    `json:"VersionNumber"`
	Description            string `json:"Description,omitempty"`
	ContentType            string `json:"ContentType"`
	Content                []byte `json:"-"`
}

type Deployment struct {
	ApplicationID          string `json:"ApplicationId"`
	EnvironmentID          string `json:"EnvironmentId"`
	DeploymentStrategyID   string `json:"DeploymentStrategyId"`
	ConfigurationProfileID string `json:"ConfigurationProfileId"`
	DeploymentNumber       int    `json:"DeploymentNumber"`
	ConfigurationVersion   string `json:"ConfigurationVersion"`
	State                  string `json:"State"`
}
