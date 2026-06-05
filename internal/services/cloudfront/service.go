package cloudfront

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/cloudfront/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

const (
	chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type CloudFrontService struct {
	distStore storage.Backend[string, *model.Distribution]
}

func NewCloudFrontService(factory *storage.Factory) (*CloudFrontService, error) {
	distStore, _ := storage.CreateAccountAware[*model.Distribution](factory, "cloudfront", "cloudfront-distributions.json", "wal")

	return &CloudFrontService{
		distStore: distStore,
	}, nil
}

func (s *CloudFrontService) CreateDistribution(ctx context.Context, config *model.DistributionConfig) (*model.Distribution, error) {
	id := s.randomID(14)
	dist := &model.Distribution{
		ID:               id,
		Arn:              fmt.Sprintf("arn:aws:cloudfront::000000000000:distribution/%s", id),
		Status:           "Deployed",
		LastModifiedTime: time.Now(),
		DomainName:       fmt.Sprintf("%s.cloudfront.net", id),
		DistributionConfig: config,
	}

	if err := s.distStore.Put(ctx, id, dist); err != nil {
		return nil, err
	}

	log.Printf("Created CloudFront distribution: %s", id)
	return dist, nil
}

func (s *CloudFrontService) GetDistribution(ctx context.Context, id string) (*model.Distribution, error) {
	dist, ok, err := s.distStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NoSuchDistribution: The specified distribution does not exist")
	}
	return dist, nil
}

func (s *CloudFrontService) ListDistributions(ctx context.Context) ([]*model.Distribution, error) {
	return s.distStore.Scan(ctx, func(k string) bool { return true })
}

func (s *CloudFrontService) randomID(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
