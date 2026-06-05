package model

type UserPool struct {
	ID                 string            `json:"id"`
	Name               string            `json:"name"`
	Arn                string            `json:"arn"`
	Status             string            `json:"status"`
	SigningSecret      string            `json:"signingSecret"`
	SigningKeyID       string            `json:"signingKeyId"`
	SigningPublicKey   string            `json:"signingPublicKey"`
	SigningPrivateKey  string            `json:"signingPrivateKey"`
	CreationDate       int64             `json:"creationDate"`
	LastModifiedDate   int64             `json:"lastModifiedDate"`
	UserPoolTags       map[string]string `json:"userPoolTags"`
	DeletionProtection string            `json:"deletionProtection"`
}

type CognitoUser struct {
	Username          string            `json:"username"`
	UserPoolID        string            `json:"userPoolId"`
	UserStatus        string            `json:"userStatus"`
	Enabled           bool              `json:"enabled"`
	Attributes        map[string]string `json:"attributes"`
	CreationDate      int64             `json:"creationDate"`
	LastModifiedDate  int64             `json:"lastModifiedDate"`
	PasswordHash      string            `json:"passwordHash"`
	TemporaryPassword bool              `json:"temporaryPassword"`
	GroupNames        []string          `json:"groupNames"`
	SrpSalt           string            `json:"srpSalt"`
	SrpVerifier       string            `json:"srpVerifier"`
}

type UserPoolClient struct {
	ClientID                        string   `json:"clientId"`
	UserPoolID                      string   `json:"userPoolId"`
	ClientName                      string   `json:"clientName"`
	ClientSecret                    string   `json:"clientSecret"`
	GenerateSecret                  bool     `json:"generateSecret"`
	AllowedOAuthFlowsUserPoolClient bool     `json:"allowedOAuthFlowsUserPoolClient"`
	AllowedOAuthFlows               []string `json:"allowedOAuthFlows"`
	AllowedOAuthScopes              []string `json:"allowedOAuthScopes"`
	CreationDate                    int64    `json:"creationDate"`
	LastModifiedDate                int64    `json:"lastModifiedDate"`
}

type CognitoGroup struct {
	GroupName        string `json:"groupName"`
	UserPoolID       string `json:"userPoolId"`
	Description      string `json:"description"`
	RoleArn          string `json:"roleArn"`
	Precedence       int    `json:"precedence"`
	CreationDate     int64  `json:"creationDate"`
	LastModifiedDate int64  `json:"lastModifiedDate"`
}
