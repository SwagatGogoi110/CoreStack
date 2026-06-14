package dns

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/dns/model"
)

type CloudDnsHandler struct {
	service *CloudDnsService
}

func NewCloudDnsHandler(service *CloudDnsService) *CloudDnsHandler {
	return &CloudDnsHandler{
		service: service,
	}
}

func (h *CloudDnsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/dns/v1/projects/")
	ctx := context.Background()

	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	if strings.Contains(path, "/managedZones") {
		if strings.HasSuffix(path, "/managedZones") {
			if r.Method == "GET" {
				h.handleListZones(w, r, ctx)
				return
			}
			if r.Method == "POST" {
				h.handleCreateZone(w, r, ctx)
				return
			}
		}

		// Zone specific
		idx := -1
		for i, p := range parts {
			if p == "managedZones" {
				idx = i
				break
			}
		}

		if idx != -1 && len(parts) > idx+1 {
			zoneName := parts[idx+1]
			
			if strings.HasSuffix(path, "/rrsets") {
				h.handleListRrsets(w, r, ctx, zoneName)
				return
			}
			if strings.HasSuffix(path, "/changes") {
				h.handleCreateChange(w, r, ctx, zoneName)
				return
			}
			
			if r.Method == "DELETE" {
				h.handleDeleteZone(w, r, ctx, zoneName)
				return
			}
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "Cloud DNS path not implemented: %s"}}`, r.URL.Path)
}

func (h *CloudDnsHandler) handleListZones(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	items, _ := h.service.ListZones(ctx)
	res := model.ManagedZonesList{
		ManagedZones: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *CloudDnsHandler) handleCreateZone(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	var zone model.ManagedZone
	json.NewDecoder(r.Body).Decode(&zone)
	created, _ := h.service.CreateZone(ctx, zone.Name, &zone)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *CloudDnsHandler) handleListRrsets(w http.ResponseWriter, r *http.Request, ctx context.Context, zoneName string) {
	items, _ := h.service.ListRrsets(ctx, zoneName)
	res := model.ResourceRecordSetsList{
		Rrsets: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *CloudDnsHandler) handleCreateChange(w http.ResponseWriter, r *http.Request, ctx context.Context, zoneName string) {
	var change model.Change
	json.NewDecoder(r.Body).Decode(&change)
	h.service.CreateChange(ctx, zoneName, change.Additions, change.Deletions)
	change.Status = "done"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(change)
}

func (h *CloudDnsHandler) handleDeleteZone(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	h.service.DeleteZone(ctx, name)
	w.WriteHeader(http.StatusNoContent)
}
