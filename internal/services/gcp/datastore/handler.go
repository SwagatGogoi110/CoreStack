package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/datastore/model"
)

type DatastoreHandler struct {
	service *DatastoreService
}

func NewDatastoreHandler(service *DatastoreService) *DatastoreHandler {
	return &DatastoreHandler{
		service: service,
	}
}

func (h *DatastoreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/")
	ctx := context.Background()

	if strings.Contains(path, ":commit") {
		project := strings.Split(strings.TrimSuffix(path, ":commit"), "/")[1]
		h.handleCommit(w, r, ctx, project)
		return
	}

	if strings.Contains(path, ":lookup") {
		h.handleLookup(w, r, ctx)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "Datastore path not implemented: %s"}}`, path)
}

func (h *DatastoreHandler) handleCommit(w http.ResponseWriter, r *http.Request, ctx context.Context, project string) {
	var req struct {
		Mutations []map[string]any `json:"mutations"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	
	results, _ := h.service.Commit(ctx, project, req.Mutations)
	res := model.CommitResponse{
		MutationResults: results,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *DatastoreHandler) handleLookup(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	var req struct {
		Keys []*model.Key `json:"keys"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	
	lookupRes, _ := h.service.Lookup(ctx, req.Keys)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lookupRes)
}
