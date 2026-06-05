package model

type HostedZone struct {
	ID              string `json:"Id"`
	Name            string `json:"Name"`
	CallerReference string `json:"CallerReference"`
	Config          *HostedZoneConfig `json:"Config,omitempty"`
	ResourceRecordSetCount int64 `json:"ResourceRecordSetCount"`
}

type HostedZoneConfig struct {
	Comment     string `json:"Comment,omitempty"`
	PrivateZone bool   `json:"PrivateZone,omitempty"`
}

type ResourceRecordSet struct {
	Name            string            `json:"Name"`
	Type            string            `json:"Type"`
	TTL             int64             `json:"TTL,omitempty"`
	ResourceRecords []*ResourceRecord `json:"ResourceRecords,omitempty"`
}

type ResourceRecord struct {
	Value string `json:"Value"`
}
