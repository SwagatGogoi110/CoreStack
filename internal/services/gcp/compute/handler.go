package compute

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/compute/model"
)

type ComputeEngineHandler struct {
	service *ComputeEngineService
}

func NewComputeEngineHandler(service *ComputeEngineService) *ComputeEngineHandler {
	return &ComputeEngineHandler{
		service: service,
	}
}

func (h *ComputeEngineHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/compute/v1/projects/")
	ctx := context.Background()

	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	project := parts[0]

	if strings.Contains(path, "/zones") {
		if len(parts) == 2 && parts[1] == "zones" {
			h.handleListZones(w, r, ctx, project)
			return
		}
		if len(parts) >= 3 {
			zone := parts[2]
			if len(parts) >= 4 && parts[3] == "instances" {
				if len(parts) == 4 {
					if r.Method == "GET" {
						h.handleListInstances(w, r, ctx, project, zone)
						return
					}
					if r.Method == "POST" {
						h.handleCreateInstance(w, r, ctx, project, zone)
						return
					}
				}
				if len(parts) == 5 {
					instance := parts[4]
					if r.Method == "GET" {
						h.handleGetInstance(w, r, ctx, project, zone, instance)
						return
					}
					if r.Method == "DELETE" {
						h.handleDeleteInstance(w, r, ctx, project, zone, instance)
						return
					}
				}
			}
		}
	}

	if strings.Contains(path, "/regions") {
		h.handleListRegions(w, r, ctx, project)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "Compute path not implemented: %s"}}`, r.URL.Path)
}

func (h *ComputeEngineHandler) handleListInstances(w http.ResponseWriter, r *http.Request, ctx context.Context, project, zone string) {
	items, _ := h.service.ListInstances(ctx, project, zone)
	res := model.InstancesList{
		Kind:  "compute#instanceList",
		Items: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *ComputeEngineHandler) handleCreateInstance(w http.ResponseWriter, r *http.Request, ctx context.Context, project, zone string) {
	var inst model.Instance
	json.NewDecoder(r.Body).Decode(&inst)
	
	created, _ := h.service.CreateInstance(ctx, project, zone, inst.Name, &inst)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *ComputeEngineHandler) handleGetInstance(w http.ResponseWriter, r *http.Request, ctx context.Context, project, zone, name string) {
	items, _ := h.service.ListInstances(ctx, project, zone)
	for _, inst := range items {
		if inst.Name == name {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(inst)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func (h *ComputeEngineHandler) handleDeleteInstance(w http.ResponseWriter, r *http.Request, ctx context.Context, project, zone, name string) {
	h.service.DeleteInstance(ctx, project, zone, name)
	w.WriteHeader(http.StatusNoContent)
}

func (h *ComputeEngineHandler) handleListZones(w http.ResponseWriter, r *http.Request, ctx context.Context, project string) {
	res := model.ZonesList{
		Kind: "compute#zoneList",
		Items: []*model.Zone{
			{Name: "us-central1-a", Status: "UP", Region: "us-central1"},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *ComputeEngineHandler) handleListRegions(w http.ResponseWriter, r *http.Request, ctx context.Context, project string) {
	res := model.RegionsList{
		Kind: "compute#regionList",
		Items: []*model.Region{
			{Name: "us-central1", Status: "UP"},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
