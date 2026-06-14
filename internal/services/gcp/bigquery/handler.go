package bigquery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/bigquery/model"
)

type BigQueryHandler struct {
	service *BigQueryService
}

func NewBigQueryHandler(service *BigQueryService) *BigQueryHandler {
	return &BigQueryHandler{
		service: service,
	}
}

func (h *BigQueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/bigquery/v2/")
	ctx := context.Background()

	// Datasets
	if strings.Contains(path, "/datasets") {
		parts := strings.Split(path, "/")
		project := parts[1]

		if len(parts) == 3 && parts[2] == "datasets" {
			if r.Method == "GET" {
				h.handleListDatasets(w, r, ctx, project)
				return
			}
			if r.Method == "POST" {
				h.handleCreateDataset(w, r, ctx, project)
				return
			}
		}
	}

	// Tables
	if strings.Contains(path, "/tables") {
		parts := strings.Split(path, "/")
		project := parts[1]
		datasetId := parts[3]

		if len(parts) == 5 && parts[4] == "tables" {
			if r.Method == "GET" {
				h.handleListTables(w, r, ctx, project, datasetId)
				return
			}
			if r.Method == "POST" {
				h.handleCreateTable(w, r, ctx, project, datasetId)
				return
			}
		}
	}

	// Queries
	if strings.Contains(path, "/queries") {
		parts := strings.Split(path, "/")
		project := parts[1]
		if r.Method == "POST" {
			h.handleQuery(w, r, ctx, project)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "BigQuery path not implemented: %s"}}`, path)
}

func (h *BigQueryHandler) handleListDatasets(w http.ResponseWriter, r *http.Request, ctx context.Context, project string) {
	datasets, _ := h.service.ListDatasets(ctx, project)
	res := model.DatasetsList{
		Kind:     "bigquery#datasetList",
		Datasets: datasets,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *BigQueryHandler) handleCreateDataset(w http.ResponseWriter, r *http.Request, ctx context.Context, project string) {
	var req model.Dataset
	json.NewDecoder(r.Body).Decode(&req)
	
	created, err := h.service.CreateDataset(ctx, project, req.DatasetReference.DatasetId, &req)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *BigQueryHandler) handleListTables(w http.ResponseWriter, r *http.Request, ctx context.Context, project, datasetId string) {
	tables, _ := h.service.ListTables(ctx, project, datasetId)
	res := model.TablesList{
		Kind:   "bigquery#tableList",
		Tables: tables,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *BigQueryHandler) handleCreateTable(w http.ResponseWriter, r *http.Request, ctx context.Context, project, datasetId string) {
	var req model.Table
	json.NewDecoder(r.Body).Decode(&req)
	
	created, err := h.service.CreateTable(ctx, project, datasetId, req.TableReference.TableId, &req)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *BigQueryHandler) handleQuery(w http.ResponseWriter, r *http.Request, ctx context.Context, project string) {
	var req struct {
		Query string `json:"query"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	
	res, _ := h.service.Query(ctx, project, req.Query)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
