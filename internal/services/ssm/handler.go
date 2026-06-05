package ssm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type SsmJsonHandler struct {
	service *SsmService
}

func NewSsmJsonHandler(service *SsmService) *SsmJsonHandler {
	return &SsmJsonHandler{
		service: service,
	}
}

func (h *SsmJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "PutParameter":
		var req struct {
			Name        string `json:"Name"`
			Value       string `json:"Value"`
			Type        string `json:"Type"`
			Overwrite   bool   `json:"Overwrite"`
			Description string `json:"Description"`
		}
		json.Unmarshal(request, &req)
		version, err := h.service.PutParameter(ctx, req.Name, req.Value, req.Type, req.Overwrite)
		if err != nil {
			return nil, err
		}
		return map[string]int64{"Version": version}, nil

	case "GetParameter":
		var req struct {
			Name string `json:"Name"`
		}
		json.Unmarshal(request, &req)
		param, err := h.service.GetParameter(ctx, req.Name)
		if err != nil {
			return nil, err
		}
		return map[string]any{"Parameter": param}, nil

	case "GetParameters":
		var req struct {
			Names []string `json:"Names"`
		}
		json.Unmarshal(request, &req)
		params, _ := h.service.GetParameters(ctx, req.Names)
		return map[string]any{"Parameters": params, "InvalidParameters": []string{}}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *SsmJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("SSM does not support Query protocol")
}
