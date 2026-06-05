package opensearch

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/opensearch/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type OpenSearchService struct {
	domainStore storage.Backend[string, *model.DomainStatus]
}

func NewOpenSearchService(factory *storage.Factory) (*OpenSearchService, error) {
	domainStore, _ := storage.CreateAccountAware[*model.DomainStatus](factory, "opensearch", "os-domains.json", "wal")

	return &OpenSearchService{
		domainStore: domainStore,
	}, nil
}

func (s *OpenSearchService) CreateDomain(ctx context.Context, name, version string) (*model.DomainStatus, error) {
	if _, ok, _ := s.domainStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("ResourceAlreadyExistsException: Domain already exists")
	}

	domain := &model.DomainStatus{
		DomainName:    name,
		DomainID:      "000000000000/" + name,
		Arn:           fmt.Sprintf("arn:aws:es:us-east-1:000000000000:domain/%s", name),
		Created:       true,
		Processing:    false,
		EngineVersion: version,
		Endpoint:      fmt.Sprintf("%s.us-east-1.es.localhost", name),
		CreatedAt:     time.Now(),
	}

	if err := s.domainStore.Put(ctx, name, domain); err != nil {
		return nil, err
	}

	log.Printf("Created OpenSearch domain: %s", name)
	return domain, nil
}

func (s *OpenSearchService) DescribeDomain(ctx context.Context, name string) (*model.DomainStatus, error) {
	domain, ok, err := s.domainStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("ResourceNotFoundException: Domain not found")
	}
	return domain, nil
}

func (s *OpenSearchService) ListDomains(ctx context.Context) ([]string, error) {
	domains, err := s.domainStore.Scan(ctx, func(k string) bool { return true })
	if err != nil {
		return nil, err
	}
	var names []string
	for _, d := range domains {
		names = append(names, d.DomainName)
	}
	return names, nil
}
