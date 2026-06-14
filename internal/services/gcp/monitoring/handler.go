package monitoring

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/monitoring/model"
)

type CloudMonitoringHandler struct {
	service *CloudMonitoringService
}

func NewCloudMonitoringHandler(service *CloudMonitoringService) *CloudMonitoringHandler {
	return &CloudMonitoringHandler{
		service: service,
	}
}

func (h *CloudMonitoringHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if strings.Contains(r.URL.Path, "/timeSeries") {
		if r.Method == "POST" {
			h.handleCreateTimeSeries(w, r, ctx)
			return
		}
		if r.Method == "GET" {
			h.handleListTimeSeries(w, r, ctx)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "Monitoring path not implemented: %s"}}`, r.URL.Path)
}

func (h *CloudMonitoringHandler) handleCreateTimeSeries(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	var req struct {
		TimeSeries []*model.TimeSeries `json:"timeSeries"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	h.service.CreateTimeSeries(ctx, req.TimeSeries)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "{}")
}

func (h *CloudMonitoringHandler) handleListTimeSeries(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	items, _ := h.service.ListTimeSeries(ctx)
	res := model.ListTimeSeriesResponse{
		TimeSeries: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
