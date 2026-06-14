package appengine

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/appengine/model"
)

type AppEngineHandler struct {
	service *AppEngineService
}

func NewAppEngineHandler(service *AppEngineService) *AppEngineHandler {
	return &AppEngineHandler{
		service: service,
	}
}

func (h *AppEngineHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/apps/")
	ctx := context.Background()

	parts := strings.Split(path, "/")
	project := parts[0]

	if len(parts) == 1 {
		if r.Method == "GET" {
			h.handleGetApplication(w, r, ctx, project)
			return
		}
	}

	if len(parts) >= 2 && parts[1] == "services" {
		if len(parts) == 2 {
			if r.Method == "GET" {
				h.handleListServices(w, r, ctx)
				return
			}
		}
		
		if len(parts) >= 4 && parts[3] == "versions" {
			if r.Method == "GET" {
				h.handleListVersions(w, r, ctx, parts[2])
				return
			}
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "App Engine path not implemented: %s"}}`, r.URL.Path)
}

func (h *AppEngineHandler) handleGetApplication(w http.ResponseWriter, r *http.Request, ctx context.Context, id string) {
	app, err := h.service.GetApplication(ctx, id)
	if err != nil {
		// Auto-create for health check if needed
		app = &model.Application{
			Id: id,
			LocationId: "us-central",
			DefaultHostname: id + ".appspot.com",
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app)
}

func (h *AppEngineHandler) handleListServices(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	items, _ := h.service.ListServices(ctx)
	res := model.ServicesList{
		Services: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *AppEngineHandler) handleListVersions(w http.ResponseWriter, r *http.Request, ctx context.Context, serviceId string) {
	items, _ := h.service.ListVersions(ctx, serviceId)
	res := model.VersionsList{
		Versions: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
