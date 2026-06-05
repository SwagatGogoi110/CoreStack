package model

type Export struct {
	ExportArn                string                    `json:"ExportArn"`
	Name                     string                    `json:"Name"`
	Description              string                    `json:"Description,omitempty"`
	DataQuery                *DataQuery                `json:"DataQuery"`
	DestinationConfigurations *DestinationConfiguration `json:"DestinationConfigurations"`
	RefreshCadence           *RefreshCadence           `json:"RefreshCadence"`
	CreatedAt                int64                     `json:"CreatedAt"`
	LastUpdatedAt            int64                     `json:"LastUpdatedAt"`
	ExportStatus             string                    `json:"ExportStatus"` // HEALTHY | UNHEALTHY
}

type DataQuery struct {
	QueryStatement string `json:"QueryStatement"`
}

type DestinationConfiguration struct {
	S3Destination *S3Destination `json:"S3Destination"`
}

type S3Destination struct {
	S3Bucket string `json:"S3Bucket"`
	S3Prefix string `json:"S3Prefix"`
	S3Region string `json:"S3Region"`
}

type RefreshCadence struct {
	Frequency string `json:"Frequency"` // SYNCHRONOUS
}
