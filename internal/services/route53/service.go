package route53

import (
	"context"
	"log"
	"math/rand"
	"sync"

	"github.com/hectorvent/cloudstack/internal/services/route53/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

const (
	chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type Route53Service struct {
	zoneStore storage.Backend[string, *model.HostedZone]
	mu        sync.RWMutex
}

func NewRoute53Service(factory *storage.Factory) (*Route53Service, error) {
	zoneStore, _ := storage.CreateAccountAware[*model.HostedZone](factory, "route53", "r53-zones.json", "wal")

	return &Route53Service{
		zoneStore: zoneStore,
	}, nil
}

func (s *Route53Service) CreateHostedZone(ctx context.Context, name, callerRef string) (*model.HostedZone, error) {
	id := "/hostedzone/" + s.randomID(14)
	zone := &model.HostedZone{
		ID:              id,
		Name:            name,
		CallerReference: callerRef,
		ResourceRecordSetCount: 2,
	}

	if err := s.zoneStore.Put(ctx, id, zone); err != nil {
		return nil, err
	}

	log.Printf("Created Route53 hosted zone: %s (%s)", name, id)
	return zone, nil
}

func (s *Route53Service) ListHostedZones(ctx context.Context) ([]*model.HostedZone, error) {
	return s.zoneStore.Scan(ctx, func(k string) bool { return true })
}

func (s *Route53Service) randomID(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
