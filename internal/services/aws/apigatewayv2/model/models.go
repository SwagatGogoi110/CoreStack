package model

type Api struct {
	ApiID                    string            `json:"apiId"`
	Name                     string            `json:"name"`
	ProtocolType             string            `json:"protocolType"` // HTTP, WEBSOCKET
	ApiEndpoint              string            `json:"apiEndpoint"`
	CreatedDate              int64             `json:"createdDate"`
	Tags                     map[string]string `json:"tags,omitempty"`
	RouteSelectionExpression string            `json:"routeSelectionExpression,omitempty"`
	Description              string            `json:"description,omitempty"`
}

type Route struct {
	RouteID                    string `json:"routeId"`
	RouteKey                   string `json:"routeKey"`
	AuthorizationType          string `json:"authorizationType"`
	Target                     string `json:"target,omitempty"`
	RouteResponseSelectionExpr string `json:"routeResponseSelectionExpression,omitempty"`
}

type Integration struct {
	IntegrationID        string `json:"integrationId"`
	IntegrationType      string `json:"integrationType"`
	IntegrationUri       string `json:"integrationUri,omitempty"`
	PayloadFormatVersion string `json:"payloadFormatVersion,omitempty"`
	IntegrationMethod    string `json:"integrationMethod,omitempty"`
}
