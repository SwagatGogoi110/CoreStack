package apprunner

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type ApprunnerJsonHandler struct {
	service *ApprunnerService
}

func NewApprunnerJsonHandler(service *ApprunnerService) *ApprunnerJsonHandler {
	return &ApprunnerJsonHandler{
		service: service,
	}
}

func (h *ApprunnerJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateService":
		var req struct {
			ServiceName string `json:"ServiceName"`
		}
		json.Unmarshal(request, &req)
		svc, err := h.service.CreateService(ctx, req.ServiceName)
		if err != nil {
			return nil, err
		}
		return map[string]any{"Service": svc}, nil

	case "DescribeService":
		var req struct {
			ServiceArn string `json:"ServiceArn"`
		}
		json.Unmarshal(request, &req)
		svc, err := h.service.DescribeService(ctx, req.ServiceArn)
		if err != nil {
			return nil, err
		}
		return map[string]any{"Service": svc}, nil

	case "ListServices":
		summaries, _ := h.service.ListServices(ctx)
		return map[string]any{"ServiceSummaryList": summaries}, nil

	case "DeleteService":
		var req struct {
			ServiceArn string `json:"ServiceArn"`
		}
		json.Unmarshal(request, &req)
		h.service.DeleteService(ctx, req.ServiceArn)
		return map[string]any{}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *ApprunnerJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("AppRunner does not support Query protocol")
}
