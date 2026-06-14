package waf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type WafJsonHandler struct {
	service *WafService
}

func NewWafJsonHandler(service *WafService) *WafJsonHandler {
	return &WafJsonHandler{
		service: service,
	}
}

func (h *WafJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateWebACL":
		var req struct {
			Name          string `json:"Name"`
			MetricName    string `json:"MetricName"`
			DefaultAction struct {
				Type string `json:"Type"`
			} `json:"DefaultAction"`
		}
		json.Unmarshal(request, &req)
		acl, err := h.service.CreateWebACL(ctx, req.Name, req.MetricName, req.DefaultAction.Type)
		if err != nil {
			return nil, err
		}
		return map[string]any{"WebACL": acl}, nil

	case "GetWebACL":
		var req struct {
			WebACLId string `json:"WebACLId"`
		}
		json.Unmarshal(request, &req)
		acl, err := h.service.GetWebACL(ctx, req.WebACLId)
		if err != nil {
			return nil, err
		}
		return map[string]any{"WebACL": acl}, nil

	case "ListWebACLs":
		acls, _ := h.service.ListWebACLs(ctx)
		return map[string]any{"WebACLs": acls}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *WafJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("WAF does not support Query protocol")
}
