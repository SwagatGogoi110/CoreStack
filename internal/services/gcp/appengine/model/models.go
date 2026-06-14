package model

import "time"

type Application struct {
	Id           string `json:"id"`
	LocationId   string `json:"locationId"`
	AuthDomain   string `json:"authDomain"`
	DefaultHostname string `json:"defaultHostname"`
}

type Service struct {
	Name   string          `json:"name"`
	Id     string          `json:"id"`
	Split  *TrafficSplit   `json:"split"`
}

type TrafficSplit struct {
	ShardBy    string            `json:"shardBy"`
	Allocations map[string]float64 `json:"allocations"`
}

type Version struct {
	Name        string    `json:"name"`
	Id          string    `json:"id"`
	Runtime     string    `json:"runtime"`
	Entrypoint  *Entrypoint `json:"entrypoint,omitempty"`
	CreateTime  time.Time `json:"createTime"`
}

type Entrypoint struct {
	Shell string `json:"shell"`
}

type ServicesList struct {
	Services      []*Service `json:"services"`
	NextPageToken string     `json:"nextPageToken,omitempty"`
}

type VersionsList struct {
	Versions      []*Version `json:"versions"`
	NextPageToken string     `json:"nextPageToken,omitempty"`
}
