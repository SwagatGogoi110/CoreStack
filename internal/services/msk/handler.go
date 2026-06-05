package msk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type MskJsonHandler struct {
	service *MskService
}

func NewMskJsonHandler(service *MskService) *MskJsonHandler {
	return &MskJsonHandler{
		service: service,
	}
}

func (h *MskJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
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
		return cluster, nil

	case "DescribeCluster":
		var req struct {
			ClusterArn string `json:"clusterArn"`
		}
		json.Unmarshal(request, &req)
		cluster, err := h.service.DescribeCluster(ctx, req.ClusterArn)
		if err != nil {
			return nil, err
		}
		return map[string]any{"clusterInfo": cluster}, nil

	case "ListClusters":
		clusters, _ := h.service.ListClusters(ctx)
		return map[string]any{"clusterInfoList": clusters}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *MskJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("MSK does not support Query protocol")
}
