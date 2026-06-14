package elbv2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type ElbV2QueryHandler struct {
	service *ElbV2Service
}

func NewElbV2QueryHandler(service *ElbV2Service) *ElbV2QueryHandler {
	return &ElbV2QueryHandler{
		service: service,
	}
}

func (h *ElbV2QueryHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("ELBv2 does not support JSON protocol")
}

func (h *ElbV2QueryHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	ctx := context.Background()

	switch action {
	case "CreateLoadBalancer":
		name := params.Get("Name")
		scheme := params.Get("Scheme")
		lbType := params.Get("Type")
		lb, err := h.service.CreateLoadBalancer(ctx, name, scheme, lbType)
		if err != nil {
			return "", err
		}
		
		b := common.NewXmlBuilder().Start("CreateLoadBalancerResponse").Start("CreateLoadBalancerResult").Start("LoadBalancers")
		b.Start("member").
			Elem("LoadBalancerArn", lb.LoadBalancerArn).
			Elem("DNSName", lb.DNSName).
			Elem("LoadBalancerName", lb.LoadBalancerName).
			End()
		b.End().End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-elb").End().End()
		return b.Build(), nil

	case "DescribeLoadBalancers":
		lbs, _ := h.service.DescribeLoadBalancers(ctx, nil)
		
		b := common.NewXmlBuilder().Start("DescribeLoadBalancersResponse").Start("DescribeLoadBalancersResult").Start("LoadBalancers")
		for _, lb := range lbs {
			b.Start("member").
				Elem("LoadBalancerArn", lb.LoadBalancerArn).
				Elem("DNSName", lb.DNSName).
				Elem("LoadBalancerName", lb.LoadBalancerName).
				End()
		}
		b.End().End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-elb").End().End()
		return b.Build(), nil

	default:
		return "", fmt.Errorf("Unknown ELBv2 action: %s", action)
	}
}
