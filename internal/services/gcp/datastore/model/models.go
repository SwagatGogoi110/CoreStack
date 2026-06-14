package model

type Entity struct {
	Key        *Key             `json:"key"`
	Properties map[string]Value `json:"properties"`
}

type Key struct {
	PartitionId *PartitionId `json:"partitionId"`
	Path        []*PathElement `json:"path"`
}

type PartitionId struct {
	ProjectId   string `json:"projectId"`
	NamespaceId string `json:"namespaceId"`
}

type PathElement struct {
	Kind string `json:"kind"`
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Value struct {
	NullValue     any            `json:"nullValue,omitempty"`
	BooleanValue  *bool          `json:"booleanValue,omitempty"`
	IntegerValue  *string        `json:"integerValue,omitempty"`
	DoubleValue   *float64       `json:"doubleValue,omitempty"`
	TimestampValue *string       `json:"timestampValue,omitempty"`
	StringValue   *string        `json:"stringValue,omitempty"`
	BlobValue     []byte         `json:"blobValue,omitempty"`
	EntityValue   *Entity        `json:"entityValue,omitempty"`
	ArrayValue    *ArrayValue    `json:"arrayValue,omitempty"`
}

type ArrayValue struct {
	Values []Value `json:"values"`
}

type LookupResponse struct {
	Found    []*EntityResult `json:"found"`
	Missing  []*EntityResult `json:"missing"`
	Deferred []*Key          `json:"deferred"`
}

type EntityResult struct {
	Entity *Entity `json:"entity"`
	Cursor string  `json:"cursor"`
}

type CommitResponse struct {
	MutationResults []*MutationResult `json:"mutationResults"`
	IndexUpdates    int               `json:"indexUpdates"`
}

type MutationResult struct {
	Key             *Key   `json:"key"`
	Version         string `json:"version"`
	ConflictDetected bool   `json:"conflictDetected"`
}
