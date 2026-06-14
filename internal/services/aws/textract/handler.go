package textract

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type TextractJsonHandler struct {
	service *TextractService
}

func NewTextractJsonHandler(service *TextractService) *TextractJsonHandler {
	return &TextractJsonHandler{
		service: service,
	}
}

func (h *TextractJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	switch action {
	case "ListAdapters":
		return map[string]any{"Adapters": []any{}}, nil
	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *TextractJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListAdapters" {
		b := common.NewXmlBuilder()
		// Namespace?
		b.Start("ListAdaptersResponse").Start("ListAdaptersResult")
		b.Start("Adapters").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
