package model

type Dataset struct {
	Kind              string            `json:"kind"`
	Id                string            `json:"id"`
	DatasetReference  *DatasetReference `json:"datasetReference"`
	FriendlyName      string            `json:"friendlyName,omitempty"`
	Description       string            `json:"description,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	CreationTime      int64             `json:"creationTime,string"`
	LastModifiedTime  int64             `json:"lastModifiedTime,string"`
	Location          string            `json:"location"`
	DefaultTableExpirationMs int64       `json:"defaultTableExpirationMs,string,omitempty"`
}

type DatasetReference struct {
	ProjectId string `json:"projectId"`
	DatasetId string `json:"datasetId"`
}

type Table struct {
	Kind           string          `json:"kind"`
	Id             string          `json:"id"`
	TableReference *TableReference `json:"tableReference"`
	Schema         *TableSchema    `json:"schema,omitempty"`
	NumBytes       int64           `json:"numBytes,string"`
	NumRows        uint64          `json:"numRows,string"`
	CreationTime   int64           `json:"creationTime,string"`
	LastModifiedTime int64         `json:"lastModifiedTime,string"`
	Type           string          `json:"type"`
	Location       string          `json:"location"`
}

type TableReference struct {
	ProjectId string `json:"projectId"`
	DatasetId string `json:"datasetId"`
	TableId   string `json:"tableId"`
}

type TableSchema struct {
	Fields []*TableField `json:"fields"`
}

type TableField struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Mode string `json:"mode,omitempty"`
}

type Job struct {
	Kind          string         `json:"kind"`
	Id            string         `json:"id"`
	JobReference  *JobReference  `json:"jobReference"`
	Configuration *JobConfiguration `json:"configuration"`
	Status        *JobStatus     `json:"status"`
	Statistics    *JobStatistics `json:"statistics,omitempty"`
}

type JobReference struct {
	ProjectId string `json:"projectId"`
	JobId     string `json:"jobId"`
	Location  string `json:"location"`
}

type JobConfiguration struct {
	Query *JobConfigurationQuery `json:"query,omitempty"`
}

type JobConfigurationQuery struct {
	Query string `json:"query"`
}

type JobStatus struct {
	State string `json:"state"`
	ErrorResult *ErrorProto `json:"errorResult,omitempty"`
}

type ErrorProto struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type JobStatistics struct {
	CreationTime int64 `json:"creationTime,string"`
	StartTime    int64 `json:"startTime,string"`
	EndTime      int64 `json:"endTime,string"`
}

type DatasetsList struct {
	Kind          string     `json:"kind"`
	Datasets      []*Dataset `json:"datasets"`
	NextPageToken string     `json:"nextPageToken,omitempty"`
}

type TablesList struct {
	Kind          string   `json:"kind"`
	Tables        []*Table `json:"tables"`
	NextPageToken string   `json:"nextPageToken,omitempty"`
}

type QueryResponse struct {
	Kind          string         `json:"kind"`
	Schema        *TableSchema   `json:"schema"`
	JobReference  *JobReference  `json:"jobReference"`
	Rows          []*TableRow    `json:"rows"`
	TotalRows     uint64         `json:"totalRows,string"`
	JobComplete   bool           `json:"jobComplete"`
}

type TableRow struct {
	F []*TableCell `json:"f"`
}

type TableCell struct {
	V any `json:"v"`
}
