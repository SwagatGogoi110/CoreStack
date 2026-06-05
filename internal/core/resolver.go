package core

import (
	"regexp"
)

var (
	credentialRegionPattern = regexp.MustCompile(`Credential=\S+/\d{8}/([^/]+)/`)
	akidPattern             = regexp.MustCompile(`Credential=([^/]+)/`)
	twelveDigits            = regexp.MustCompile(`^\d{12}$`)
)

// Resolver handles region and account ID resolution.
type Resolver struct {
	DefaultRegion    string
	DefaultAccountID string
}

// NewResolver creates a new Resolver.
func NewResolver(defaultRegion, defaultAccountID string) *Resolver {
	if defaultRegion == "" {
		defaultRegion = "us-east-1"
	}
	if defaultAccountID == "" {
		defaultAccountID = "000000000000"
	}
	return &Resolver{
		DefaultRegion:    defaultRegion,
		DefaultAccountID: defaultAccountID,
	}
}

// ResolveRegion extracts the region from the Authorization header.
func (res *Resolver) ResolveRegion(auth string) string {
	if auth == "" {
		return res.DefaultRegion
	}
	matches := credentialRegionPattern.FindStringSubmatch(auth)
	if len(matches) > 1 {
		return matches[1]
	}
	return res.DefaultRegion
}

// ResolveAccount extracts the account ID from the Authorization header.
func (res *Resolver) ResolveAccount(auth string) string {
	if auth == "" {
		return res.DefaultAccountID
	}
	matches := akidPattern.FindStringSubmatch(auth)
	if len(matches) > 1 {
		akid := matches[1]
		if twelveDigits.MatchString(akid) {
			return akid
		}
	}
	return res.DefaultAccountID
}
