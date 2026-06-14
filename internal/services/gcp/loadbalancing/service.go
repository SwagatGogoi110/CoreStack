package loadbalancing

import (
	"context"
	"fmt"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/gcp/loadbalancing/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type LoadBalancingService struct {
	frStore   storage.Backend[string, *model.ForwardingRule]
	proxyStore storage.Backend[string, *model.TargetHttpProxy]
	urlStore   storage.Backend[string, *model.UrlMap]
	beStore    storage.Backend[string, *model.BackendService]
}

func NewLoadBalancingService(factory *storage.Factory) (*LoadBalancingService, error) {
	frStore, _ := storage.CreateAccountAware[*model.ForwardingRule](factory, "lb", "lb-fr.json", "wal")
	proxyStore, _ := storage.CreateAccountAware[*model.TargetHttpProxy](factory, "lb", "lb-proxy.json", "wal")
	urlStore, _ := storage.CreateAccountAware[*model.UrlMap](factory, "lb", "lb-url.json", "wal")
	beStore, _ := storage.CreateAccountAware[*model.BackendService](factory, "lb", "lb-be.json", "wal")

	return &LoadBalancingService{
		frStore:   frStore,
		proxyStore: proxyStore,
		urlStore:   urlStore,
		beStore:    beStore,
	}, nil
}

// Forwarding Rules

func (s *LoadBalancingService) CreateForwardingRule(ctx context.Context, name string, fr *model.ForwardingRule) (*model.ForwardingRule, error) {
	fr.Name = name
	fr.CreationTimestamp = time.Now()
	fr.SelfLink = fmt.Sprintf("https://www.googleapis.com/compute/v1/global/forwardingRules/%s", name)
	s.frStore.Put(ctx, name, fr)
	return fr, nil
}

func (s *LoadBalancingService) ListForwardingRules(ctx context.Context) ([]*model.ForwardingRule, error) {
	return s.frStore.Scan(ctx, func(k string) bool { return true })
}

// Target Proxies

func (s *LoadBalancingService) CreateTargetHttpProxy(ctx context.Context, name string, proxy *model.TargetHttpProxy) (*model.TargetHttpProxy, error) {
	proxy.Name = name
	proxy.CreationTimestamp = time.Now()
	proxy.SelfLink = fmt.Sprintf("https://www.googleapis.com/compute/v1/global/targetHttpProxies/%s", name)
	s.proxyStore.Put(ctx, name, proxy)
	return proxy, nil
}

func (s *LoadBalancingService) ListTargetHttpProxies(ctx context.Context) ([]*model.TargetHttpProxy, error) {
	return s.proxyStore.Scan(ctx, func(k string) bool { return true })
}

// URL Maps

func (s *LoadBalancingService) CreateUrlMap(ctx context.Context, name string, um *model.UrlMap) (*model.UrlMap, error) {
	um.Name = name
	um.CreationTimestamp = time.Now()
	um.SelfLink = fmt.Sprintf("https://www.googleapis.com/compute/v1/global/urlMaps/%s", name)
	s.urlStore.Put(ctx, name, um)
	return um, nil
}

func (s *LoadBalancingService) ListUrlMaps(ctx context.Context) ([]*model.UrlMap, error) {
	return s.urlStore.Scan(ctx, func(k string) bool { return true })
}
func (s *LoadBalancingService) CreateBackendService(ctx context.Context, name string, be *model.BackendService) (*model.BackendService, error) {
	be.Name = name
	be.CreationTimestamp = time.Now()
	be.SelfLink = fmt.Sprintf("https://www.googleapis.com/compute/v1/global/backendServices/%s", name)
	s.beStore.Put(ctx, name, be)
	return be, nil
}

func (s *LoadBalancingService) ListBackendServices(ctx context.Context) ([]*model.BackendService, error) {
	return s.beStore.Scan(ctx, func(k string) bool { return true })
}
