package resourcegroupstagging

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
	"github.com/hectorvent/cloudstack/internal/services/resourcegroupstagging/model"
)

type ResourceGroupsTaggingJsonHandler struct {
	service *ResourceGroupsTaggingService
}

func NewResourceGroupsTaggingJsonHandler(service *ResourceGroupsTaggingService) *ResourceGroupsTaggingJsonHandler {
	return &ResourceGroupsTaggingJsonHandler{
		service: service,
	}
}

func (h *ResourceGroupsTaggingJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "TagResources":
		var req struct {
			ResourceARNList []string          `json:"ResourceARNList"`
			Tags            map[string]string `json:"Tags"`
		}
		json.Unmarshal(request, &req)
		if err := h.service.TagResources(ctx, req.ResourceARNList, req.Tags); err != nil {
			return nil, err
		}
		return map[string]any{"FailedResourcesMap": map[string]any{}}, nil

	case "GetResources":
		mappings, _ := h.service.GetResources(ctx)
		res := make([]map[string]any, 0)
		for _, m := range mappings {
			tags := make([]model.Tag, 0)
			for k, v := range m.Tags {
				tags = append(tags, model.Tag{Key: k, Value: v})
			}
			res = append(res, map[string]any{
				"ResourceARN": m.ResourceArn,
				"Tags":        tags,
			})
		}
		return map[string]any{"ResourceTagMappingList": res}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *ResourceGroupsTaggingJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("ResourceGroupsTagging does not support Query protocol")
}
