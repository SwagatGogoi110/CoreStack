package appconfig

import (
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
	// Handle with and without prefix
	if action == "ListApplications" || action == "AppConfig.ListApplications" {
		return map[string]any{"Items": []any{}}, nil
	}
	return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
}

func (h *AppConfigJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListApplications" || action == "AppConfig.ListApplications" {
		b := common.NewXmlBuilder()
		b.Start("ListApplicationsResponse").Start("ListApplicationsResult")
		b.Start("Items").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
