package eventbridge

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type EventBridgeHandler struct {
	service *EventBridgeService
}

func NewEventBridgeHandler(service *EventBridgeService) *EventBridgeHandler {
	return &EventBridgeHandler{
		service: service,
	}
}

func (h *EventBridgeHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	switch action {
	case "ListEventBuses":
		return map[string]any{"EventBuses": []any{}}, nil
	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *EventBridgeHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListEventBuses" {
		b := common.NewXmlBuilder()
		// Namespace?
		b.Start("ListEventBusesResponse").Start("ListEventBusesResult")
		b.Start("EventBuses").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
