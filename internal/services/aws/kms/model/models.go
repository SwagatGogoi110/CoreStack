package model

type KmsKey struct {
	KeyID                 string            `json:"KeyId"`
	Arn                   string            `json:"Arn"`
	Description           string            `json:"Description,omitempty"`
	Enabled               bool              `json:"Enabled"`
	KeyState              string            `json:"KeyState"` // Enabled, Disabled, PendingDeletion
	KeyUsage              string            `json:"KeyUsage"` // ENCRYPT_DECRYPT
	CustomerMasterKeySpec string            `json:"CustomerMasterKeySpec"`
	CreationDate          int64             `json:"CreationDate"`
	DeletionDate          int64             `json:"DeletionDate,omitempty"`
	Tags                  map[string]string `json:"-"`
}

type KmsAlias struct {
	AliasName string `json:"AliasName"`
	AliasArn  string `json:"AliasArn"`
	TargetKeyID string `json:"TargetKeyId"`
}
