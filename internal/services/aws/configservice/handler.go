package configservice

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
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
	// Handle with and without prefix
	if action == "DescribeConfigRules" || action == "StarlingDoveService.DescribeConfigRules" {
		return map[string]any{"ConfigRules": []any{}}, nil
	}
	return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
}

func (h *ConfigServiceJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "DescribeConfigRules" || action == "StarlingDoveService.DescribeConfigRules" {
		b := common.NewXmlBuilder()
		b.Start("DescribeConfigRulesResponse").Start("DescribeConfigRulesResult")
		b.Start("ConfigRules").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
