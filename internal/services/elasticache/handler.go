package elasticache

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type ElastiCacheQueryHandler struct {
	service *ElastiCacheService
}

func NewElastiCacheQueryHandler(service *ElastiCacheService) *ElastiCacheQueryHandler {
	return &ElastiCacheQueryHandler{
		service: service,
	}
}

func (h *ElastiCacheQueryHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("ElastiCache does not support JSON protocol")
}

func (h *ElastiCacheQueryHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	ctx := context.Background()

	switch action {
	case "CreateReplicationGroup":
		id := params.Get("ReplicationGroupId")
		desc := params.Get("ReplicationGroupDescription")
		g, err := h.service.CreateReplicationGroup(ctx, id, desc)
		if err != nil {
			return "", err
		}
		
		b := common.NewXmlBuilder().Start("CreateReplicationGroupResponse").Start("CreateReplicationGroupResult")
		b.Start("ReplicationGroup").
			Elem("ReplicationGroupId", g.ReplicationGroupID).
			Elem("Status", g.Status).
			End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-ec").End().End()
		return b.Build(), nil

	case "DescribeReplicationGroups":
		id := params.Get("ReplicationGroupId")
		groups, _ := h.service.DescribeReplicationGroups(ctx, id)
		
		b := common.NewXmlBuilder().Start("DescribeReplicationGroupsResponse").Start("DescribeReplicationGroupsResult").Start("ReplicationGroups")
		for _, g := range groups {
			b.Start("ReplicationGroup").
				Elem("ReplicationGroupId", g.ReplicationGroupID).
				Elem("Status", g.Status).
				Start("PrimaryEndpoint").
				Elem("Address", g.PrimaryEndpoint.Address).
				Elem("Port", fmt.Sprintf("%d", g.PrimaryEndpoint.Port)).
				End().
				End()
		}
		b.End().End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-ec").End().End()
		return b.Build(), nil

	default:
		return "", fmt.Errorf("Unknown ElastiCache action: %s", action)
	}
}
