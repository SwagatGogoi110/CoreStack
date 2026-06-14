package appsync

import (
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
	// Handle with and without prefix
	if action == "ListGraphqlApis" || action == "AppSync.ListGraphqlApis" {
		return map[string]any{"graphqlApis": []any{}}, nil
	}
	return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
}

func (h *AppSyncJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListGraphqlApis" || action == "AppSync.ListGraphqlApis" {
		b := common.NewXmlBuilder()
		b.Start("ListGraphqlApisResponse").Start("ListGraphqlApisResult")
		b.Start("graphqlApis").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
