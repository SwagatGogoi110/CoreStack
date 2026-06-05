package common

type ServiceProtocol string

const (
	ProtocolJSON     ServiceProtocol = "JSON"
	ProtocolCBOR     ServiceProtocol = "CBOR"
	ProtocolQUERY    ServiceProtocol = "QUERY"
	ProtocolRESTJSON ServiceProtocol = "REST_JSON"
	ProtocolRESTXML  ServiceProtocol = "REST_XML"
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
	byCredentialScope   map[string]*ServiceDescriptor
}

type targetPrefixEntry struct {
	prefix     string
	descriptor *ServiceDescriptor
}

func NewCatalog(descriptors []*ServiceDescriptor) *Catalog {
	var targetPrefixes []targetPrefixEntry
	byCredentialScope := make(map[string]*ServiceDescriptor)
	
	for _, d := range descriptors {
		for _, prefix := range d.TargetPrefixes {
			targetPrefixes = append(targetPrefixes, targetPrefixEntry{
				prefix:     prefix,
				descriptor: d,
			})
		}
		for _, scope := range d.CredentialScopes {
			byCredentialScope[scope] = d
		}
	}
	
	return &Catalog{
		descriptors:         descriptors,
		targetPrefixes:      targetPrefixes,
		byCredentialScope:   byCredentialScope,
	}
}

func (c *Catalog) ByCredentialScope(scope string) *ServiceDescriptor {
	return c.byCredentialScope[scope]
}


func (c *Catalog) MatchTarget(target string) *TargetMatch {
	for _, entry := range c.targetPrefixes {
		if len(target) >= len(entry.prefix) && target[:len(entry.prefix)] == entry.prefix {
			return &TargetMatch{
				Descriptor: entry.descriptor,
				Prefix:     entry.prefix,
				Action:     target[len(entry.prefix):],
			}
		}
	}
	return nil
}
