package bigtable

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/bigtable/model"
)

type BigtableHandler struct {
	service *BigtableService
}

func NewBigtableHandler(service *BigtableService) *BigtableHandler {
	return &BigtableHandler{
		service: service,
	}
}

func (h *BigtableHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v2/")
	ctx := context.Background()

	if strings.Contains(path, "/instances") {
		parts := strings.Split(path, "/")
		project := parts[1]

		if len(parts) == 3 && parts[2] == "instances" {
			if r.Method == "GET" {
				h.handleListInstances(w, r, ctx, project)
				return
			}
			if r.Method == "POST" {
				h.handleCreateInstance(w, r, ctx, project)
				return
			}
		}

		if len(parts) >= 5 && parts[4] == "tables" {
			if len(parts) == 5 {
				if r.Method == "GET" {
					h.handleListTables(w, r, ctx, path)
					return
				}
				if r.Method == "POST" {
					h.handleCreateTable(w, r, ctx, path)
					return
				}
			}
			
			if strings.Contains(path, ":mutateRows") {
				h.handleMutateRows(w, r, ctx)
				return
			}
			if strings.Contains(path, ":readRows") {
				h.handleReadRows(w, r, ctx)
				return
			}
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "Bigtable path not implemented: %s"}}`, path)
}

func (h *BigtableHandler) handleListInstances(w http.ResponseWriter, r *http.Request, ctx context.Context, project string) {
	items, _ := h.service.ListInstances(ctx, project)
	res := model.InstancesList{
		Instances: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *BigtableHandler) handleCreateInstance(w http.ResponseWriter, r *http.Request, ctx context.Context, project string) {
	var req struct {
		InstanceId string          `json:"instanceId"`
		Instance   *model.Instance `json:"instance"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	name := fmt.Sprintf("projects/%s/instances/%s", project, req.InstanceId)
	created, _ := h.service.CreateInstance(ctx, name, req.Instance)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *BigtableHandler) handleListTables(w http.ResponseWriter, r *http.Request, ctx context.Context, instanceName string) {
	items, _ := h.service.ListTables(ctx, instanceName)
	res := model.TablesList{
		Tables: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *BigtableHandler) handleCreateTable(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	var table model.Table
	json.NewDecoder(r.Body).Decode(&table)
	
	tableId := r.URL.Query().Get("tableId")
	if tableId == "" {
		// Fallback for some clients
	}
	
	name := parent + "/tables/" + tableId
	created, _ := h.service.CreateTable(ctx, name, &table)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *BigtableHandler) handleMutateRows(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"entries": []any{}})
}

func (h *BigtableHandler) handleReadRows(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]any{}) // Array of rows
}
