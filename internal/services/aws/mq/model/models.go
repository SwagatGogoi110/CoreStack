package model

type BrokerSummary struct {
	BrokerArn     string `json:"BrokerArn"`
	BrokerId      string `json:"BrokerId"`
	BrokerName    string `json:"BrokerName"`
	BrokerState   string `json:"BrokerState"`
	DeploymentMode string `json:"DeploymentMode"`
	EngineType    string `json:"EngineType"`
}

type BrokerInstance struct {
	ConsoleURL string   `json:"ConsoleURL"`
	Endpoints  []string `json:"Endpoints"`
	IpAddress  string   `json:"IpAddress"`
}
