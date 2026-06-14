package model

import "time"

type ManagedZone struct {
	Name        string    `json:"name"`
	DnsName     string    `json:"dnsName"`
	Description string    `json:"description,omitempty"`
	Id          uint64    `json:"id,string"`
	CreationTime time.Time `json:"creationTime"`
	NameServers []string  `json:"nameServers,omitempty"`
	Visibility  string    `json:"visibility"`
}

type ResourceRecordSet struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Ttl     int      `json:"ttl"`
	Rrdatas []string `json:"rrdatas"`
}

type Change struct {
	Id            string               `json:"id"`
	StartTime     time.Time            `json:"startTime"`
	Status        string               `json:"status"`
	Additions     []*ResourceRecordSet `json:"additions,omitempty"`
	Deletions     []*ResourceRecordSet `json:"deletions,omitempty"`
}

type ManagedZonesList struct {
	ManagedZones  []*ManagedZone `json:"managedZones"`
	NextPageToken string         `json:"nextPageToken,omitempty"`
}

type ResourceRecordSetsList struct {
	Rrsets        []*ResourceRecordSet `json:"rrsets"`
	NextPageToken string               `json:"nextPageToken,omitempty"`
}
