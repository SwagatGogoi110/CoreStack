package model

type RestApi struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	CreatedDate int64             `json:"createdDate"`
	Tags        map[string]string `json:"tags,omitempty"`
}

type Resource struct {
	ID       string `json:"id"`
	ParentID string `json:"parentId,omitempty"`
	PathPart string `json:"pathPart,omitempty"`
	Path     string `json:"path"`
}

type Method struct {
	HTTPMethod      string `json:"httpMethod"`
	AuthType        string `json:"authorizationType"`
	Integration     *Integration `json:"methodIntegration,omitempty"`
}

type Integration struct {
	Type       string `json:"type"`
	HTTPMethod string `json:"httpMethod"`
	Uri        string `json:"uri"`
}

type Deployment struct {
	ID          string `json:"id"`
	Description string `json:"description,omitempty"`
	CreatedDate int64  `json:"createdDate"`
}

type Stage struct {
	StageName    string `json:"stageName"`
	DeploymentID string `json:"deploymentId"`
}
