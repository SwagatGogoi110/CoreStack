package resourcegroupstagging

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
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
	switch action {
	case "GetResources":
		return map[string]any{"ResourceTagMappingList": []any{}}, nil
	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *ResourceGroupsTaggingJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "GetResources" {
		b := common.NewXmlBuilder()
		// Namespace?
		b.Start("GetResourcesResponse").Start("GetResourcesResult")
		b.Start("ResourceTagMappingList").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
