package model

import (
	"time"
)

type Identity struct {
	Identity           string    `json:"Identity"`
	IdentityType       string    `json:"IdentityType"` // EmailAddress, Domain
	VerificationStatus string    `json:"VerificationStatus"` // Success, Pending
	CreatedAt          time.Time `json:"CreatedAt"`
}

type SentEmail struct {
	MessageId    string    `json:"MessageId"`
	Source       string    `json:"Source"`
	Destinations []string  `json:"Destinations"`
	Subject      string    `json:"Subject"`
	Timestamp    time.Time `json:"Timestamp"`
}
