package bedrockruntime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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
	return nil, fmt.Errorf("BedrockRuntime does not support standard JSON protocol dispatcher")
}

func (h *BedrockRuntimeHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("BedrockRuntime does not support Query protocol")
}

func (h *BedrockRuntimeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/model/")
	
	if strings.HasSuffix(path, "/converse") {
		modelID := strings.TrimSuffix(path, "/converse")
		h.handleConverse(w, r, modelID)
		return
	}
	
	if strings.HasSuffix(path, "/invoke") {
		modelID := strings.TrimSuffix(path, "/invoke")
		h.handleInvoke(w, r, modelID)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func (h *BedrockRuntimeHandler) handleConverse(w http.ResponseWriter, r *http.Request, modelID string) {
	res := h.service.BuildConverseResponse(modelID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *BedrockRuntimeHandler) handleInvoke(w http.ResponseWriter, r *http.Request, modelID string) {
	res := h.service.BuildInvokeModelResponse(modelID)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
