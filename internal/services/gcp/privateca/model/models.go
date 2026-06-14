package model

import "time"

type CaPool struct {
	Name string `json:"name"`
	Tier string `json:"tier"`
}

type CertificateAuthority struct {
	Name            string          `json:"name"`
	Type            string          `json:"type"`
	State           string          `json:"state"`
	Config          *CaConfig       `json:"config"`
	CreateTime      time.Time       `json:"createTime"`
}

type CaConfig struct {
	SubjectConfig *SubjectConfig `json:"subjectConfig"`
}

type SubjectConfig struct {
	Subject *Subject `json:"subject"`
}

type Subject struct {
	CommonName   string `json:"commonName"`
	Organization string `json:"organization"`
}

type Certificate struct {
	Name            string          `json:"name"`
	PemCertificate  string          `json:"pemCertificate"`
	CreateTime      time.Time       `json:"createTime"`
}

type CaPoolsList struct {
	CaPools       []*CaPool `json:"caPools"`
	NextPageToken string    `json:"nextPageToken,omitempty"`
}

type CertificateAuthoritiesList struct {
	CertificateAuthorities []*CertificateAuthority `json:"certificateAuthorities"`
}
