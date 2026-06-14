package cloudbuild

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/cloudbuild/model"
)

type CloudBuildHandler struct {
	service *CloudBuildService
}

func NewCloudBuildHandler(service *CloudBuildService) *CloudBuildHandler {
	return &CloudBuildHandler{
		service: service,
	}
}

func (h *CloudBuildHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/")
	ctx := context.Background()

	if strings.Contains(path, "/builds") {
		parts := strings.Split(path, "/")
		project := parts[1]
		if r.Method == "GET" {
			h.handleListBuilds(w, r, ctx, project)
			return
		}
		if r.Method == "POST" {
			h.handleCreateBuild(w, r, ctx, project)
			return
		}
	}

	if strings.Contains(path, "/triggers") {
		parts := strings.Split(path, "/")
		project := parts[1]
		if r.Method == "GET" {
			h.handleListTriggers(w, r, ctx, project)
			return
		}
		if r.Method == "POST" {
			h.handleCreateTrigger(w, r, ctx, project)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "Cloud Build path not implemented: %s"}}`, path)
}

func (h *CloudBuildHandler) handleListBuilds(w http.ResponseWriter, r *http.Request, ctx context.Context, project string) {
	items, _ := h.service.ListBuilds(ctx, project)
	res := model.BuildsList{
		Builds: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *CloudBuildHandler) handleCreateBuild(w http.ResponseWriter, r *http.Request, ctx context.Context, project string) {
	var build model.Build
	json.NewDecoder(r.Body).Decode(&build)
	created, _ := h.service.CreateBuild(ctx, project, &build)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *CloudBuildHandler) handleListTriggers(w http.ResponseWriter, r *http.Request, ctx context.Context, project string) {
	items, _ := h.service.ListTriggers(ctx, project)
	res := model.TriggersList{
		Triggers: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *CloudBuildHandler) handleCreateTrigger(w http.ResponseWriter, r *http.Request, ctx context.Context, project string) {
	var trigger model.BuildTrigger
	json.NewDecoder(r.Body).Decode(&trigger)
	created, _ := h.service.CreateTrigger(ctx, project, &trigger)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}
