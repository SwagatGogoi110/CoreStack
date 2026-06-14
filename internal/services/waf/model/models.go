package model

type WebACL struct {
	DefaultAction *WafAction `json:"DefaultAction"`
	MetricName    string     `json:"MetricName"`
	Name          string     `json:"Name"`
	WebACLId      string     `json:"WebACLId"`
}

type WafAction struct {
	Type string `json:"Type"`
}

type Rule struct {
	MetricName string `json:"MetricName"`
	Name       string `json:"Name"`
	RuleId     string `json:"RuleId"`
}
