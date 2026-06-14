package model

type Organization struct {
	Id           string `json:"Id"`
	Arn          string `json:"Arn"`
	FeatureSet   string `json:"FeatureSet"`
	MasterAccountId string `json:"MasterAccountId"`
	MasterAccountArn string `json:"MasterAccountArn"`
	MasterAccountEmail string `json:"MasterAccountEmail"`
}

type Account struct {
	Id     string `json:"Id"`
	Arn    string `json:"Arn"`
	Email  string `json:"Email"`
	Name   string `json:"Name"`
	Status string `json:"Status"`
	JoinedMethod string `json:"JoinedMethod"`
}
