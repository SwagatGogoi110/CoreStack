package model

type Trail struct {
	Name                       string `json:"Name"`
	S3BucketName               string `json:"S3BucketName"`
	IncludeGlobalServiceEvents bool   `json:"IncludeGlobalServiceEvents"`
	IsMultiRegionTrail         bool   `json:"IsMultiRegionTrail"`
	HomeRegion                 string `json:"HomeRegion"`
	TrailARN                   string `json:"TrailARN"`
	LogFileValidationEnabled   bool   `json:"LogFileValidationEnabled"`
}
