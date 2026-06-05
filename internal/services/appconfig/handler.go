package appconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type AppConfigJsonHandler struct {
	service *AppConfigService
}

func NewAppConfigJsonHandler(service *AppConfigService) *AppConfigJsonHandler {
	return &AppConfigJsonHandler{
		service: service,
	}
}

func (h *AppConfigJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateApplication":
		var req struct {
			Name        string `json:"Name"`
			Description string `json:"Description"`
		}
		json.Unmarshal(request, &req)
		app, err := h.service.CreateApplication(ctx, req.Name, req.Description)
		if err != nil {
			return nil, err
		}
		return app, nil

	case "GetApplication":
		var req struct {
			ApplicationID string `json:"ApplicationId"`
		}
		json.Unmarshal(request, &req)
		app, err := h.service.GetApplication(ctx, req.ApplicationID)
		if err != nil {
			return nil, err
		}
		return app, nil

	case "ListApplications":
		apps, _ := h.service.ListApplications(ctx)
		return map[string]any{"Items": apps}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *AppConfigJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("AppConfig does not support Query protocol")
}
