package model

type GraphqlApi struct {
	ApiID             string            `json:"apiId"`
	Name              string            `json:"name"`
	Arn               string            `json:"arn"`
	AuthenticationType string            `json:"authenticationType"`
	Uris              map[string]string `json:"uris"`
	Tags              map[string]string `json:"tags,omitempty"`
}

type DataSource struct {
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	Type           string `json:"type"`
	ServiceRoleArn string `json:"serviceRoleArn,omitempty"`
}

type Resolver struct {
	TypeName                 string `json:"typeName"`
	FieldName                string `json:"fieldName"`
	DataSourceName           string `json:"dataSourceName,omitempty"`
	RequestMappingTemplate   string `json:"requestMappingTemplate,omitempty"`
	ResponseMappingTemplate  string `json:"responseMappingTemplate,omitempty"`
	Kind                     string `json:"kind"`
}
