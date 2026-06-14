package model

import "time"

type Service struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Uri         string            `json:"uri"`
	Template    *RevisionTemplate `json:"template"`
	Traffic     []*TrafficTarget  `json:"traffic"`
	CreateTime  time.Time         `json:"createTime"`
	UpdateTime  time.Time         `json:"updateTime"`
	Etag        string            `json:"etag"`
}

type RevisionTemplate struct {
	Revision   string      `json:"revision"`
	Containers []*Container `json:"containers"`
}

type Container struct {
	Image string `json:"image"`
}

type TrafficTarget struct {
	Type     string `json:"type"`
	Revision string `json:"revision"`
	Percent  int    `json:"percent"`
}

type ServicesList struct {
	Services      []*Service `json:"services"`
	NextPageToken string     `json:"nextPageToken,omitempty"`
}
