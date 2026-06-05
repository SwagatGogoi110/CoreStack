package ce

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type CostExplorerJsonHandler struct {
	service *CostExplorerService
}

func NewCostExplorerJsonHandler(service *CostExplorerService) *CostExplorerJsonHandler {
	return &CostExplorerJsonHandler{
		service: service,
	}
}

func (h *CostExplorerJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	switch action {
	case "GetCostAndUsage", "GetCostAndUsageWithResources":
		return h.service.GetCostAndUsage(), nil
	case "GetDimensionValues":
		return h.service.GetDimensionValues(), nil
	case "GetTags":
		return map[string]any{"Tags": []string{}, "ReturnSize": 0, "TotalSize": 0}, nil
	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *CostExplorerJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("CostExplorer does not support Query protocol")
}
