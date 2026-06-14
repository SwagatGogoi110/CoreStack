package bedrockruntime

import (
	"net/http"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type BedrockRuntimeHandler struct {
	service *BedrockRuntimeService
}

func NewBedrockRuntimeHandler(service *BedrockRuntimeService) *BedrockRuntimeHandler {
	return &BedrockRuntimeHandler{
		service: service,
	}
}

func (h *BedrockRuntimeHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	switch action {
	case "ListCustomModels", "ListAsyncInvokes":
		return map[string]any{"modelSummaries": []any{}, "asyncInvocationSummaries": []any{}}, nil
	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *BedrockRuntimeHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListCustomModels" || action == "ListCustomModels" {
		b := common.NewXmlBuilder()
		b.Start("ListCustomModelsResponse").Start("ListCustomModelsResult")
		b.Start("modelSummaries").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}

func (h *BedrockRuntimeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
