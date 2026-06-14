package cloudtrail

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type CloudtrailJsonHandler struct {
	service *CloudtrailService
}

func NewCloudtrailJsonHandler(service *CloudtrailService) *CloudtrailJsonHandler {
	return &CloudtrailJsonHandler{
		service: service,
	}
}

func (h *CloudtrailJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	// Handle with and without prefix
	if action == "DescribeTrails" || action == "CloudTrail_20131101.DescribeTrails" {
		return map[string]any{"trailList": []any{}}, nil
	}
	return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
}

func (h *CloudtrailJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "DescribeTrails" || action == "CloudTrail_20131101.DescribeTrails" {
		b := common.NewXmlBuilder()
		b.Start("DescribeTrailsResponse").Start("DescribeTrailsResult")
		b.Start("trailList").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
