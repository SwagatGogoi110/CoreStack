package appsync

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type AppSyncJsonHandler struct {
	service *AppSyncService
}

func NewAppSyncJsonHandler(service *AppSyncService) *AppSyncJsonHandler {
	return &AppSyncJsonHandler{
		service: service,
	}
}

func (h *AppSyncJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateGraphqlApi":
		var req struct {
			Name               string `json:"name"`
			AuthenticationType string `json:"authenticationType"`
		}
		json.Unmarshal(request, &req)
		api, err := h.service.CreateGraphqlApi(ctx, req.Name, req.AuthenticationType)
		if err != nil {
			return nil, err
		}
		return map[string]any{"graphqlApi": api}, nil

	case "GetGraphqlApi":
		var req struct {
			ApiID string `json:"apiId"`
		}
		json.Unmarshal(request, &req)
		api, err := h.service.GetGraphqlApi(ctx, req.ApiID)
		if err != nil {
			return nil, err
		}
		return map[string]any{"graphqlApi": api}, nil

	case "ListGraphqlApis":
		apis, _ := h.service.ListGraphqlApis(ctx)
		return map[string]any{"graphqlApis": apis}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *AppSyncJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("AppSync does not support Query protocol")
}
