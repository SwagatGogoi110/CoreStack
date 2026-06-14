package model

type Operation struct {
	Name     string         `json:"name"`
	Metadata map[string]any `json:"metadata,omitempty"`
	Done     bool           `json:"done"`
	Error    *Status        `json:"error,omitempty"`
	Response map[string]any `json:"response,omitempty"`
}

type Status struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Details []map[string]any `json:"details,omitempty"`
}

type OperationsList struct {
	Operations    []*Operation `json:"operations"`
	NextPageToken string       `json:"nextPageToken,omitempty"`
}
