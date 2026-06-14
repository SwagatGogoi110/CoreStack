package glue

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type GlueJsonHandler struct {
	service *GlueService
}

func NewGlueJsonHandler(service *GlueService) *GlueJsonHandler {
	return &GlueJsonHandler{
		service: service,
	}
}

func (h *GlueJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	switch action {
	case "GetDatabases":
		return map[string]any{"DatabaseList": []any{}}, nil
	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *GlueJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "GetDatabases" {
		b := common.NewXmlBuilder()
		// Namespace?
		b.Start("GetDatabasesResponse").Start("GetDatabasesResult")
		b.Start("DatabaseList").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
