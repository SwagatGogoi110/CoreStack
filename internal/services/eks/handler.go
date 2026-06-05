package eks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type EksJsonHandler struct {
	service *EksService
}

func NewEksJsonHandler(service *EksService) *EksJsonHandler {
	return &EksJsonHandler{
		service: service,
	}
}

func (h *EksJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateCluster":
		var req struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}
		json.Unmarshal(request, &req)
		cluster, err := h.service.CreateCluster(ctx, req.Name, req.Version)
		if err != nil {
			return nil, err
		}
		return map[string]any{"cluster": cluster}, nil

	case "DescribeCluster":
		var req struct {
			Name string `json:"name"`
		}
		json.Unmarshal(request, &req)
		cluster, err := h.service.DescribeCluster(ctx, req.Name)
		if err != nil {
			return nil, err
		}
		return map[string]any{"cluster": cluster}, nil

	case "ListClusters":
		names, _ := h.service.ListClusters(ctx)
		return map[string]any{"clusters": names}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *EksJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("EKS does not support Query protocol")
}
