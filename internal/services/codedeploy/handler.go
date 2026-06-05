package codedeploy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type CodeDeployJsonHandler struct {
	service *CodeDeployService
}

func NewCodeDeployJsonHandler(service *CodeDeployService) *CodeDeployJsonHandler {
	return &CodeDeployJsonHandler{
		service: service,
	}
}

func (h *CodeDeployJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateApplication":
		var req struct {
			ApplicationName string `json:"applicationName"`
			ComputePlatform string `json:"computePlatform"`
		}
		json.Unmarshal(request, &req)
		app, err := h.service.CreateApplication(ctx, req.ApplicationName, req.ComputePlatform)
		if err != nil {
			return nil, err
		}
		return map[string]any{"applicationId": app.ApplicationID}, nil

	case "GetApplication":
		var req struct {
			ApplicationName string `json:"applicationName"`
		}
		json.Unmarshal(request, &req)
		app, err := h.service.GetApplication(ctx, req.ApplicationName)
		if err != nil {
			return nil, err
		}
		return map[string]any{"application": app}, nil

	case "ListApplications":
		names, _ := h.service.ListApplications(ctx)
		return map[string]any{"applications": names}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *CodeDeployJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("CodeDeploy does not support Query protocol")
}
