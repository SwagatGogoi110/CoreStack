package model

type BackupVault struct {
	BackupVaultName        string            `json:"BackupVaultName"`
	BackupVaultArn         string            `json:"BackupVaultArn"`
	EncryptionKeyArn       string            `json:"EncryptionKeyArn,omitempty"`
	CreationDate           int64             `json:"CreationDate"`
	NumberOfRecoveryPoints int64             `json:"NumberOfRecoveryPoints"`
	Tags                   map[string]string `json:"Tags,omitempty"`
}

type BackupPlan struct {
	BackupPlanID   string        `json:"BackupPlanId"`
	BackupPlanArn  string        `json:"BackupPlanArn"`
	BackupPlanName string        `json:"BackupPlanName"`
	CreationDate   int64         `json:"CreationDate"`
	VersionID      string        `json:"VersionId"`
	Rules          []*BackupRule `json:"Rules"`
}

type BackupRule struct {
	RuleName  string `json:"RuleName"`
	TargetVaultName string `json:"TargetBackupVaultName"`
	RuleID    string `json:"RuleId,omitempty"`
}
