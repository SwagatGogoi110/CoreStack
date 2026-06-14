package cloudwatch

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type CloudWatchQueryHandler struct {
	service *CloudWatchService
}

func NewCloudWatchQueryHandler(service *CloudWatchService) *CloudWatchQueryHandler {
	return &CloudWatchQueryHandler{
		service: service,
	}
}

func (h *CloudWatchQueryHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	// Handle with and without prefix
	if action == "ListMetrics" || action == "ListMetrics" {
		return map[string]any{"Metrics": []any{}}, nil
	}
	return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
}

func (h *CloudWatchQueryHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListMetrics" || action == "ListMetrics" {
		b := common.NewXmlBuilder()
		b.Start("ListMetricsResponse").Start("ListMetricsResult")
		b.Start("Metrics").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
