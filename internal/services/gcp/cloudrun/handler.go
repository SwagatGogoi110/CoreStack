package cloudrun

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/cloudrun/model"
)

type CloudRunHandler struct {
	service *CloudRunService
}

func NewCloudRunHandler(service *CloudRunService) *CloudRunHandler {
	return &CloudRunHandler{
		service: service,
	}
}

func (h *CloudRunHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v2/")
	ctx := context.Background()

	if strings.Contains(path, "/services") {
		if strings.HasSuffix(path, "/services") {
			if r.Method == "GET" {
				h.handleListServices(w, r, ctx, path)
				return
			}
			if r.Method == "POST" {
				h.handleCreateService(w, r, ctx, path)
				return
			}
		}

		if r.Method == "GET" {
			h.handleGetService(w, r, ctx, path)
			return
		}
		if r.Method == "DELETE" {
			h.handleDeleteService(w, r, ctx, path)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "GCP Cloud Run path not implemented: %s"}}`, path)
}

func (h *CloudRunHandler) handleListServices(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	services, _ := h.service.ListServices(ctx, parent)
	res := model.ServicesList{
		Services: services,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *CloudRunHandler) handleCreateService(w http.ResponseWriter, r *http.Request, ctx context.Context, path string) {
	var svc model.Service
	json.NewDecoder(r.Body).Decode(&svc)
	
	serviceId := r.URL.Query().Get("serviceId")
	name := path + "/" + serviceId

	created, err := h.service.CreateService(ctx, name, &svc)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *CloudRunHandler) handleGetService(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	svc, err := h.service.GetService(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(svc)
}

func (h *CloudRunHandler) handleDeleteService(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	h.service.DeleteService(ctx, name)
	w.WriteHeader(http.StatusNoContent)
}
