package model

import "time"

type ForwardingRule struct {
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	IPAddress   string    `json:"IPAddress"`
	IPProtocol  string    `json:"IPProtocol"`
	PortRange   string    `json:"portRange"`
	Target      string    `json:"target"`
	SelfLink    string    `json:"selfLink"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
}

type TargetHttpProxy struct {
	Name        string    `json:"name"`
	UrlMap      string    `json:"urlMap"`
	SelfLink    string    `json:"selfLink"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
}

type UrlMap struct {
	Name           string    `json:"name"`
	DefaultService string    `json:"defaultService"`
	SelfLink       string    `json:"selfLink"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
}

type BackendService struct {
	Name        string    `json:"name"`
	Protocol    string    `json:"protocol"`
	TimeoutSec  int       `json:"timeoutSec"`
	SelfLink    string    `json:"selfLink"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
}

type ForwardingRulesList struct {
	Items []*ForwardingRule `json:"items"`
}

type TargetHttpProxiesList struct {
	Items []*TargetHttpProxy `json:"items"`
}

type UrlMapsList struct {
	Items []*UrlMap `json:"items"`
}

type BackendServicesList struct {
	Items []*BackendService `json:"items"`
}
