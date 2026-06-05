package pricing

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type PricingJsonHandler struct {
	service *PricingService
}

func NewPricingJsonHandler(service *PricingService) *PricingJsonHandler {
	return &PricingJsonHandler{
		service: service,
	}
}

func (h *PricingJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	switch action {
	case "DescribeServices":
		return h.service.DescribeServices(), nil
	case "GetProducts":
		return h.service.GetProducts(), nil
	case "GetAttributeValues":
		return map[string]any{"AttributeValues": []any{}}, nil
	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *PricingJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("Pricing does not support Query protocol")
}
