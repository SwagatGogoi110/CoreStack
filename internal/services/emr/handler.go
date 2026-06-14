package emr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type EmrJsonHandler struct {
	service *EmrService
}

func NewEmrJsonHandler(service *EmrService) *EmrJsonHandler {
	return &EmrJsonHandler{
		service: service,
	}
}

func (h *EmrJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "RunJobFlow":
		var req struct {
			Name         string `json:"Name"`
			ReleaseLabel string `json:"ReleaseLabel"`
		}
		json.Unmarshal(request, &req)
		cluster, err := h.service.RunJobFlow(ctx, req.Name, req.ReleaseLabel)
		if err != nil {
			return nil, err
		}
		return map[string]string{"JobFlowId": cluster.Id}, nil

	case "DescribeCluster":
		var req struct {
			ClusterId string `json:"ClusterId"`
		}
		json.Unmarshal(request, &req)
		cluster, err := h.service.DescribeCluster(ctx, req.ClusterId)
		if err != nil {
			return nil, err
		}
		return map[string]any{"Cluster": cluster}, nil

	case "ListClusters":
		clusters, _ := h.service.ListClusters(ctx)
		return map[string]any{"Clusters": clusters}, nil

	case "TerminateJobFlows":
		var req struct {
			JobFlowIds []string `json:"JobFlowIds"`
		}
		json.Unmarshal(request, &req)
		h.service.TerminateJobFlows(ctx, req.JobFlowIds)
		return map[string]any{}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *EmrJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("EMR does not support Query protocol")
}
