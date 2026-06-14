package model

import "time"

type Bucket struct {
	Kind              string            `json:"kind"`
	Id                string            `json:"id"`
	SelfLink          string            `json:"selfLink"`
	Name              string            `json:"name"`
	ProjectNumber     string            `json:"projectNumber,omitempty"`
	Metageneration    string            `json:"metageneration"`
	Location          string            `json:"location"`
	StorageClass      string            `json:"storageClass"`
	TimeCreated       time.Time         `json:"timeCreated"`
	Updated           time.Time         `json:"updated"`
	Etag              string            `json:"etag"`
	ProjectId         string            `json:"projectId,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	Versioning        map[string]any    `json:"versioning,omitempty"`
	Lifecycle         map[string]any    `json:"lifecycle,omitempty"`
	Cors              []map[string]any  `json:"cors,omitempty"`
	RetentionPolicy   map[string]any    `json:"retentionPolicy,omitempty"`
	DefaultEventBasedHold bool          `json:"defaultEventBasedHold,omitempty"`
}

type ObjectMeta struct {
	Kind               string            `json:"kind"`
	Id                 string            `json:"id"`
	SelfLink           string            `json:"selfLink"`
	Name               string            `json:"name"`
	Bucket             string            `json:"bucket"`
	Generation         string            `json:"generation"`
	Metageneration     string            `json:"metageneration"`
	ContentType        string            `json:"contentType"`
	StorageClass       string            `json:"storageClass"`
	Size               string            `json:"size"`
	TimeCreated        time.Time         `json:"timeCreated"`
	Updated            time.Time         `json:"updated"`
	Crc32c             string            `json:"crc32c"`
	Md5Hash            string            `json:"md5Hash"`
	MediaLink          string            `json:"mediaLink"`
	Etag               string            `json:"etag"`
	ContentDisposition string            `json:"contentDisposition,omitempty"`
	ContentEncoding    string            `json:"contentEncoding,omitempty"`
	ContentLanguage    string            `json:"contentLanguage,omitempty"`
	Metadata           map[string]string `json:"metadata,omitempty"`
	CustomerEncryption map[string]string `json:"customerEncryption,omitempty"`
}

type ObjectsList struct {
	Kind          string        `json:"kind"`
	Items         []*ObjectMeta `json:"items,omitempty"`
	NextPageToken string        `json:"nextPageToken,omitempty"`
	Prefixes      []string      `json:"prefixes,omitempty"`
}

type BucketsList struct {
	Kind          string    `json:"kind"`
	Items         []*Bucket `json:"items,omitempty"`
	NextPageToken string    `json:"nextPageToken,omitempty"`
}
