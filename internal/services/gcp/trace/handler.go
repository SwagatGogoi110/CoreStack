package trace

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/trace/model"
)

type CloudTraceHandler struct {
	service *CloudTraceService
}

func NewCloudTraceHandler(service *CloudTraceService) *CloudTraceHandler {
	return &CloudTraceHandler{
		service: service,
	}
}

func (h *CloudTraceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if strings.Contains(r.URL.Path, "/traces:batchWrite") {
		h.handleBatchWrite(w, r, ctx)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "Trace path not implemented: %s"}}`, r.URL.Path)
}

func (h *CloudTraceHandler) handleBatchWrite(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	var req model.BatchWriteSpansRequest
	json.NewDecoder(r.Body).Decode(&req)
	h.service.WriteSpans(ctx, req.Spans)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "{}")
}
