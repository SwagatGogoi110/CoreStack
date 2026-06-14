package model

import (
	"time"
)

type CertificateStatus string

const (
	StatusPendingValidation CertificateStatus = "PENDING_VALIDATION"
	StatusIssued            CertificateStatus = "ISSUED"
	StatusInactive          CertificateStatus = "INACTIVE"
	StatusExpired           CertificateStatus = "EXPIRED"
	StatusValidationTimedOut CertificateStatus = "VALIDATION_TIMED_OUT"
	StatusRevoked           CertificateStatus = "REVOKED"
	StatusFailed            CertificateStatus = "FAILED"
)

type CertificateType string

const (
	TypeAmazonIssued CertificateType = "AMAZON_ISSUED"
	TypePrivate      CertificateType = "PRIVATE"
	TypeImported     CertificateType = "IMPORTED"
)

type KeyAlgorithm string

const (
	KeyAlgorithmRSA2048 KeyAlgorithm = "RSA_2048"
	KeyAlgorithmRSA4096 KeyAlgorithm = "RSA_4096"
	KeyAlgorithmEC_PRIME256V1 KeyAlgorithm = "EC_prime256v1"
	KeyAlgorithmEC_SECP384R1  KeyAlgorithm = "EC_secp384r1"
	KeyAlgorithmEC_SECP521R1  KeyAlgorithm = "EC_secp521r1"
)

type Certificate struct {
	Arn                       string               `json:"CertificateArn"`
	DomainName                string               `json:"DomainName"`
	SubjectAlternativeNames   []string             `json:"SubjectAlternativeNames,omitempty"`
	Status                    CertificateStatus    `json:"Status"`
	Type                      CertificateType      `json:"Type"`
	CreatedAt                 time.Time            `json:"CreatedAt"`
	IssuedAt                  *time.Time           `json:"IssuedAt,omitempty"`
	ImportedAt                *time.Time           `json:"ImportedAt,omitempty"`
	NotBefore                 *time.Time           `json:"NotBefore,omitempty"`
	NotAfter                  *time.Time           `json:"NotAfter,omitempty"`
	Serial                    string               `json:"Serial,omitempty"`
	Subject                   string               `json:"Subject,omitempty"`
	Issuer                    string               `json:"Issuer,omitempty"`
	KeyAlgorithm              KeyAlgorithm         `json:"KeyAlgorithm"`
	SignatureAlgorithm        string               `json:"SignatureAlgorithm,omitempty"`
	InUseBy                   []string             `json:"InUseBy,omitempty"`
	Tags                      map[string]string    `json:"-"` // Tags are usually separate in list
	CertificateBody           string               `json:"-"`
	PrivateKey                string               `json:"-"`
	CertificateChain          string               `json:"-"`
}

type ListResult struct {
	CertificateSummaryList []*CertificateSummary `json:"CertificateSummaryList"`
	NextToken              string               `json:"NextToken,omitempty"`
}

type CertificateSummary struct {
	CertificateArn string `json:"CertificateArn"`
	DomainName     string `json:"DomainName"`
}
