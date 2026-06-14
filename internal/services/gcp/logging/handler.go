package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/logging/model"
)

type CloudLoggingHandler struct {
	service *CloudLoggingService
}

func NewCloudLoggingHandler(service *CloudLoggingService) *CloudLoggingHandler {
	return &CloudLoggingHandler{
		service: service,
	}
}

func (h *CloudLoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if strings.Contains(r.URL.Path, "/entries:write") {
		h.handleWriteEntries(w, r, ctx)
		return
	}
	if (strings.Contains(r.URL.Path, "/entries") || strings.Contains(r.URL.Path, "/logs")) && r.Method == "GET" {
		h.handleListEntries(w, r, ctx)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "Logging path not implemented: %s"}}`, r.URL.Path)
}

func (h *CloudLoggingHandler) handleWriteEntries(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	var req model.WriteLogEntriesRequest
	json.NewDecoder(r.Body).Decode(&req)
	h.service.WriteEntries(ctx, req.Entries)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "{}")
}

func (h *CloudLoggingHandler) handleListEntries(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	items, _ := h.service.ListEntries(ctx)
	res := model.ListLogEntriesResponse{
		Entries: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
