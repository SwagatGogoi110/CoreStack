package synthetics

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type SyntheticsJsonHandler struct {
	service *SyntheticsService
}

func NewSyntheticsJsonHandler(service *SyntheticsService) *SyntheticsJsonHandler {
	return &SyntheticsJsonHandler{
		service: service,
	}
}

func (h *SyntheticsJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	// Handle with and without prefix
	if action == "DescribeCanaries" || action == "Synthetics.DescribeCanaries" {
		return map[string]any{"Canaries": []any{}}, nil
	}
	return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
}

func (h *SyntheticsJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "DescribeCanaries" || action == "Synthetics.DescribeCanaries" {
		b := common.NewXmlBuilder()
		b.Start("DescribeCanariesResponse").Start("DescribeCanariesResult")
		b.Start("Canaries").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
