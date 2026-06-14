package dns

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/gcp/dns/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CloudDnsService struct {
	zoneStore storage.Backend[string, *model.ManagedZone]
	rrStore   storage.Backend[string, *model.ResourceRecordSet]
}

func NewCloudDnsService(factory *storage.Factory) (*CloudDnsService, error) {
	zoneStore, _ := storage.CreateAccountAware[*model.ManagedZone](factory, "dns", "dns-zones.json", "wal")
	rrStore, _ := storage.CreateAccountAware[*model.ResourceRecordSet](factory, "dns", "dns-rrsets.json", "wal")

	return &CloudDnsService{
		zoneStore: zoneStore,
		rrStore:   rrStore,
	}, nil
}

// Managed Zones

func (s *CloudDnsService) CreateZone(ctx context.Context, name string, zone *model.ManagedZone) (*model.ManagedZone, error) {
	if _, ok, _ := s.zoneStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("AlreadyExists: Zone %s already exists", name)
	}

	zone.Name = name
	zone.Id = rand.Uint64()
	zone.CreationTime = time.Now()
	zone.NameServers = []string{"ns-cloud-a1.googledomains.com."}

	if err := s.zoneStore.Put(ctx, name, zone); err != nil {
		return nil, err
	}
	return zone, nil
}

func (s *CloudDnsService) ListZones(ctx context.Context) ([]*model.ManagedZone, error) {
	return s.zoneStore.Scan(ctx, func(k string) bool { return true })
}

// Resource Record Sets

func (s *CloudDnsService) ListRrsets(ctx context.Context, zoneName string) ([]*model.ResourceRecordSet, error) {
	return s.rrStore.Scan(ctx, func(k string) bool {
		// Mock logic: return all for now or filter by zone in key
		return true
	})
}

func (s *CloudDnsService) CreateChange(ctx context.Context, zoneName string, additions, deletions []*model.ResourceRecordSet) error {
	for _, rr := range additions {
		s.rrStore.Put(ctx, zoneName+":"+rr.Name+":"+rr.Type, rr)
	}
	for _, rr := range deletions {
		s.rrStore.Delete(ctx, zoneName+":"+rr.Name+":"+rr.Type)
	}
	return nil
}

func (s *CloudDnsService) DeleteZone(ctx context.Context, name string) error {
	return s.zoneStore.Delete(ctx, name)
}
