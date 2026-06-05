package ecs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type EcsJsonHandler struct {
	service *EcsService
}

func NewEcsJsonHandler(service *EcsService) *EcsJsonHandler {
	return &EcsJsonHandler{
		service: service,
	}
}

func (h *EcsJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateCluster":
		var req struct {
			ClusterName string `json:"clusterName"`
		}
		json.Unmarshal(request, &req)
		cluster, err := h.service.CreateCluster(ctx, req.ClusterName)
		if err != nil {
			return nil, err
		}
		return map[string]any{"cluster": cluster}, nil

	case "DescribeClusters":
		var req struct {
			Clusters []string `json:"clusters"`
		}
		json.Unmarshal(request, &req)
		clusters, _ := h.service.DescribeClusters(ctx, req.Clusters)
		return map[string]any{"clusters": clusters, "failures": []any{}}, nil

	case "ListClusters":
		clusters, _ := h.service.DescribeClusters(ctx, nil)
		var arns []string
		for _, c := range clusters {
			arns = append(arns, c.ClusterArn)
		}
		return map[string]any{"clusterArns": arns}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *EcsJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("ECS does not support Query protocol")
}
