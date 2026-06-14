package model

type ServiceAccount struct {
	Name         string `json:"name"`
	ProjectId    string `json:"projectId"`
	UniqueId     string `json:"uniqueId"`
	Email        string `json:"email"`
	DisplayName  string `json:"displayName"`
	Description  string `json:"description"`
	Oauth2ClientId string `json:"oauth2ClientId"`
	Disabled     bool   `json:"disabled"`
}

type ServiceAccountsList struct {
	Accounts      []*ServiceAccount `json:"accounts"`
	NextPageToken string            `json:"nextPageToken,omitempty"`
}

type ServiceAccountKey struct {
	Name            string `json:"name"`
	PrivateKeyType  string `json:"privateKeyType"`
	KeyAlgorithm    string `json:"keyAlgorithm"`
	PrivateKeyData  []byte `json:"privateKeyData"`
	PublicKeyData   []byte `json:"publicKeyData"`
	ValidAfterTime  string `json:"validAfterTime"`
	ValidBeforeTime string `json:"validBeforeTime"`
}
