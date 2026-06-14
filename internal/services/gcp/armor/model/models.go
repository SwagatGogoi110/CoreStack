package model

import "time"

type SecurityPolicy struct {
	Name        string                `json:"name"`
	Description string                `json:"description,omitempty"`
	Rules       []*SecurityPolicyRule `json:"rules,omitempty"`
	CreateTime  time.Time             `json:"createTime"`
	SelfLink    string                `json:"selfLink"`
}

type SecurityPolicyRule struct {
	Priority    int            `json:"priority"`
	Action      string         `json:"action"`
	Preview     bool           `json:"preview"`
	Match       *Match         `json:"match"`
	Description string         `json:"description,omitempty"`
}

type Match struct {
	VersionedExpr string         `json:"versionedExpr"`
	Config        map[string]any `json:"config"`
}

type SecurityPoliciesList struct {
	SecurityPolicies []*SecurityPolicy `json:"items"`
}
