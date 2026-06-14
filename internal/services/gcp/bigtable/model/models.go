package model

type Instance struct {
	Name        string            `json:"name"`
	DisplayName string            `json:"displayName"`
	State       string            `json:"state"`
	Type        string            `json:"type"`
	Labels      map[string]string `json:"labels,omitempty"`
}

type Table struct {
	Name           string                    `json:"name"`
	ColumnFamilies map[string]*ColumnFamily  `json:"columnFamilies,omitempty"`
	Granularity    string                    `json:"granularity,omitempty"`
}

type ColumnFamily struct {
	GcRule map[string]any `json:"gcRule,omitempty"`
}

type InstancesList struct {
	Instances     []*Instance `json:"instances"`
	NextPageToken string      `json:"nextPageToken,omitempty"`
}

type TablesList struct {
	Tables        []*Table `json:"tables"`
	NextPageToken string   `json:"nextPageToken,omitempty"`
}

type ReadRowsResponse struct {
	Chunks []*CellChunk `json:"chunks"`
}

type CellChunk struct {
	RowKey    []byte `json:"rowKey"`
	FamilyName string `json:"familyName"`
	Qualifier  []byte `json:"qualifier"`
	Value      []byte `json:"value"`
}
