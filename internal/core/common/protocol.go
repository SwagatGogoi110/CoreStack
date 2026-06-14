package common

import (
	"strings"
)

type ServiceProtocol string

const (
	ProtocolJSON     ServiceProtocol = "JSON"
	ProtocolCBOR     ServiceProtocol = "CBOR"
	ProtocolQUERY    ServiceProtocol = "QUERY"
	ProtocolRESTJSON ServiceProtocol = "REST-JSON"
	ProtocolRESTXML  ServiceProtocol = "REST-XML"
)

type ServiceDescriptor struct {
	ExternalKey      string
	ConfigKey        string
	Enabled          bool
	TargetPrefixes   []string
	CredentialScopes []string
	DefaultProtocol  ServiceProtocol
}

type TargetMatch struct {
	Descriptor *ServiceDescriptor
	Prefix     string
	Action     string
}

type Catalog struct {
	descriptors         []*ServiceDescriptor
	targetPrefixes      []targetPrefixEntry
	byCredentialScope   map[string][]*ServiceDescriptor
}

type targetPrefixEntry struct {
	prefix     string
	descriptor *ServiceDescriptor
}

func NewCatalog(descriptors []*ServiceDescriptor) *Catalog {
	var targetPrefixes []targetPrefixEntry
	byCredentialScope := make(map[string][]*ServiceDescriptor)
	
	for _, d := range descriptors {
		for _, prefix := range d.TargetPrefixes {
			targetPrefixes = append(targetPrefixes, targetPrefixEntry{
				prefix:     prefix,
				descriptor: d,
			})
		}
		for _, scope := range d.CredentialScopes {
			byCredentialScope[scope] = append(byCredentialScope[scope], d)
		}
	}
	
	return &Catalog{
		descriptors:         descriptors,
		targetPrefixes:      targetPrefixes,
		byCredentialScope:   byCredentialScope,
	}
}

func (c *Catalog) ByCredentialScope(scope string) []*ServiceDescriptor {
	return c.byCredentialScope[scope]
}


func (c *Catalog) MatchTarget(target string) *TargetMatch {
	// Try exact prefix match first
	for _, entry := range c.targetPrefixes {
		if strings.HasPrefix(target, entry.prefix) {
			return &TargetMatch{
				Descriptor: entry.descriptor,
				Prefix:     entry.prefix,
				Action:     target[len(entry.prefix):],
			}
		}
	}
	
	// Fallback: match by action name after the last dot
	if strings.Contains(target, ".") {
		parts := strings.Split(target, ".")
		action := parts[len(parts)-1]
		
		// Look for any service that has this target prefix as part of the full target
		// or just match the action to any service that supports it?
		// For health check, let's try to match the descriptor from the first entry
		// that has a matching prefix in the full target string.
		for _, entry := range c.targetPrefixes {
			if strings.Contains(target, entry.prefix) {
				return &TargetMatch{
					Descriptor: entry.descriptor,
					Prefix:     entry.prefix,
					Action:     action,
				}
			}
		}
	}
	
	return nil
}
