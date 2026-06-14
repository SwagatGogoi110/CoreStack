package athena

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type AthenaJsonHandler struct {
	service *AthenaService
}

func NewAthenaJsonHandler(service *AthenaService) *AthenaJsonHandler {
	return &AthenaJsonHandler{
		service: service,
	}
}

func (h *AthenaJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	switch action {
	case "ListWorkGroups":
		return map[string]any{"WorkGroups": []any{}}, nil
	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *AthenaJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListWorkGroups" {
		b := common.NewXmlBuilder()
		// Namespace?
		b.Start("ListWorkGroupsResponse").Start("ListWorkGroupsResult")
		b.Start("WorkGroups").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
