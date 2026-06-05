package model

type ConfigRule struct {
	ConfigRuleName  string `json:"ConfigRuleName"`
	ConfigRuleArn   string `json:"ConfigRuleArn"`
	ConfigRuleID    string `json:"ConfigRuleId"`
	ConfigRuleState string `json:"ConfigRuleState"` // ACTIVE, DELETING
}

type ConfigurationRecorder struct {
	Name    string `json:"name"`
	RoleArn string `json:"roleARN"`
}

type DeliveryChannel struct {
	Name         string `json:"name"`
	S3BucketName string `json:"s3BucketName"`
}
