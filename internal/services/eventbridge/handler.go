package eventbridge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
	"github.com/hectorvent/cloudstack/internal/services/eventbridge/model"
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
	ctx := context.Background()

	switch action {
	case "CreateEventBus":
		var req struct {
			Name string `json:"Name"`
		}
		json.Unmarshal(request, &req)
		bus, err := h.service.CreateEventBus(ctx, req.Name)
		if err != nil {
			return nil, err
		}
		return map[string]string{"EventBusArn": bus.Arn}, nil

	case "PutRule":
		var req struct {
			Name               string `json:"Name"`
			EventBusName       string `json:"EventBusName"`
			EventPattern       string `json:"EventPattern"`
			ScheduleExpression string `json:"ScheduleExpression"`
		}
		json.Unmarshal(request, &req)
		rule, err := h.service.PutRule(ctx, req.Name, req.EventBusName, req.EventPattern, req.ScheduleExpression)
		if err != nil {
			return nil, err
		}
		return map[string]string{"RuleArn": rule.Arn}, nil

	case "PutTargets":
		var req struct {
			Rule         string          `json:"Rule"`
			EventBusName string          `json:"EventBusName"`
			Targets      []*model.Target `json:"Targets"`
		}
		json.Unmarshal(request, &req)
		if err := h.service.PutTargets(ctx, req.Rule, req.EventBusName, req.Targets); err != nil {
			return nil, err
		}
		return map[string]any{"FailedEntryCount": 0, "FailedEntries": []any{}}, nil

	default:
		return nil, fmt.Errorf("Unknown EventBridge action: %s", action)
	}
}

func (h *EventBridgeHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("EventBridge does not support Query protocol")
}
