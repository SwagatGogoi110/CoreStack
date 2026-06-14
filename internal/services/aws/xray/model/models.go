package model

type Segment struct {
	Id      string  `json:"Id"`
	Document string `json:"Document"`
}

type TraceSummary struct {
	Id           string   `json:"Id"`
	Duration     float64  `json:"Duration"`
	HasFault     bool     `json:"HasFault"`
	HasError     bool     `json:"HasError"`
	Http         *Http    `json:"Http,omitempty"`
}

type Http struct {
	HttpURL    string `json:"HttpURL"`
	HttpStatus int    `json:"HttpStatus"`
}
