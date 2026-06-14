package pipes

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type PipesJsonHandler struct {
	service *PipesService
}

func NewPipesJsonHandler(service *PipesService) *PipesJsonHandler {
	return &PipesJsonHandler{
		service: service,
	}
}

func (h *PipesJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	// Handle with and without prefix
	if action == "ListPipes" || action == "PipeService.ListPipes" {
		return map[string]any{"Pipes": []any{}}, nil
	}
	return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
}

func (h *PipesJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListPipes" || action == "PipeService.ListPipes" {
		b := common.NewXmlBuilder()
		b.Start("ListPipesResponse").Start("ListPipesResult")
		b.Start("Pipes").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
