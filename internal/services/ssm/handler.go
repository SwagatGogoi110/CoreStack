package ssm

import (
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
	switch action {
	case "ListCommands":
		return map[string]any{"Commands": []any{}}, nil
	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *SsmJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListCommands" {
		b := common.NewXmlBuilder()
		// Namespace?
		b.Start("ListCommandsResponse").Start("ListCommandsResult")
		b.Start("Commands").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
