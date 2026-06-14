package xray

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type XrayJsonHandler struct {
	service *XrayService
}

func NewXrayJsonHandler(service *XrayService) *XrayJsonHandler {
	return &XrayJsonHandler{
		service: service,
	}
}

func (h *XrayJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	// Handle with and without prefix
	if action == "GetTraceSummaries" || action == "AWSXRay.GetTraceSummaries" {
		return map[string]any{"TraceSummaries": []any{}}, nil
	}
	return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
}

func (h *XrayJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "GetTraceSummaries" || action == "AWSXRay.GetTraceSummaries" {
		b := common.NewXmlBuilder()
		b.Start("GetTraceSummariesResponse").Start("GetTraceSummariesResult")
		b.Start("TraceSummaries").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
