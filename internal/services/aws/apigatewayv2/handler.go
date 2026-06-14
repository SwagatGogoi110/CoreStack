package apigatewayv2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type ApiGatewayV2JsonHandler struct {
	service *ApiGatewayV2Service
}

func NewApiGatewayV2JsonHandler(service *ApiGatewayV2Service) *ApiGatewayV2JsonHandler {
	return &ApiGatewayV2JsonHandler{
		service: service,
	}
}

func (h *ApiGatewayV2JsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateApi":
		var req struct {
			Name         string `json:"name"`
			ProtocolType string `json:"protocolType"`
			Description  string `json:"description"`
		}
		json.Unmarshal(request, &req)
		api, err := h.service.CreateApi(ctx, req.Name, req.ProtocolType, req.Description)
		if err != nil {
			return nil, err
		}
		return api, nil

	case "GetApi":
		var req struct {
			ApiID string `json:"apiId"`
		}
		json.Unmarshal(request, &req)
		api, err := h.service.GetApi(ctx, req.ApiID)
		if err != nil {
			return nil, err
		}
		return api, nil

	case "GetApis":
		apis, _ := h.service.GetApis(ctx)
		return map[string]any{"Items": apis}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *ApiGatewayV2JsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("ApiGatewayV2 does not support Query protocol")
}
