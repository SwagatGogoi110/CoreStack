package elbv2

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/elbv2/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

const (
	hexChars = "0123456789abcdef"
)

type ElbV2Service struct {
	lbStore storage.Backend[string, *model.LoadBalancer]
	mu      sync.RWMutex
}

func NewElbV2Service(factory *storage.Factory) (*ElbV2Service, error) {
	lbStore, _ := storage.CreateAccountAware[*model.LoadBalancer](factory, "elbv2", "elb-lbs.json", "wal")

	return &ElbV2Service{
		lbStore: lbStore,
	}, nil
}

func (s *ElbV2Service) CreateLoadBalancer(ctx context.Context, name, scheme, lbType string) (*model.LoadBalancer, error) {
	id := s.randomHex(16)
	arn := fmt.Sprintf("arn:aws:elasticloadbalancing:us-east-1:000000000000:loadbalancer/%s/%s/%s", lbType, name, id)
	
	lb := &model.LoadBalancer{
		LoadBalancerArn:  arn,
		LoadBalancerName: name,
		DNSName:          fmt.Sprintf("%s-%s.elb.localhost", name, id),
		CreatedTime:      time.Now(),
		Scheme:           scheme,
		Type:             lbType,
		State:            "active",
	}

	if err := s.lbStore.Put(ctx, arn, lb); err != nil {
		return nil, err
	}

	log.Printf("Created ELBv2 load balancer: %s", name)
	return lb, nil
}

func (s *ElbV2Service) DescribeLoadBalancers(ctx context.Context, arns []string) ([]*model.LoadBalancer, error) {
	return s.lbStore.Scan(ctx, func(k string) bool {
		if len(arns) == 0 {
			return true
		}
		for _, a := range arns {
			if a == k {
				return true
			}
		}
		return false
	})
}

func (s *ElbV2Service) randomHex(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = hexChars[rand.Intn(len(hexChars))]
	}
	return string(b)
}
