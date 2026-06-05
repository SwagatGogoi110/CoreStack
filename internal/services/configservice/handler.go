package configservice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
	"github.com/hectorvent/cloudstack/internal/services/configservice/model"
)

type ConfigServiceJsonHandler struct {
	service *ConfigService
}

func NewConfigServiceJsonHandler(service *ConfigService) *ConfigServiceJsonHandler {
	return &ConfigServiceJsonHandler{
		service: service,
	}
}

func (h *ConfigServiceJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "PutConfigRule":
		var req struct {
			ConfigRule *model.ConfigRule `json:"ConfigRule"`
		}
		json.Unmarshal(request, &req)
		rule, err := h.service.PutConfigRule(ctx, req.ConfigRule)
		if err != nil {
			return nil, err
		}
		return rule, nil

	case "DescribeConfigRules":
		var req struct {
			ConfigRuleNames []string `json:"ConfigRuleNames"`
		}
		json.Unmarshal(request, &req)
		rules, _ := h.service.DescribeConfigRules(ctx, req.ConfigRuleNames)
		return map[string]any{"ConfigRules": rules}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *ConfigServiceJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("ConfigService does not support Query protocol")
}
