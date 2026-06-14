package spanner

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/spanner/model"
)

type SpannerHandler struct {
	service *SpannerService
}

func NewSpannerHandler(service *SpannerService) *SpannerHandler {
	return &SpannerHandler{
		service: service,
	}
}

func (h *SpannerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/")
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

		if len(parts) >= 5 && parts[4] == "databases" {
			if len(parts) == 5 {
				if r.Method == "GET" {
					h.handleListDatabases(w, r, ctx, path)
					return
				}
				if r.Method == "POST" {
					h.handleCreateDatabase(w, r, ctx, path)
					return
				}
			}

			if len(parts) >= 7 && parts[6] == "sessions" {
				if len(parts) == 7 && r.Method == "POST" {
					h.handleCreateSession(w, r, ctx, path)
					return
				}
				if strings.Contains(path, ":executeSql") {
					h.handleExecuteSql(w, r, ctx)
					return
				}
			}
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "Spanner path not implemented: %s"}}`, path)
}

func (h *SpannerHandler) handleListInstances(w http.ResponseWriter, r *http.Request, ctx context.Context, project string) {
	items, _ := h.service.ListInstances(ctx, project)
	res := model.InstancesList{
		Instances: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *SpannerHandler) handleCreateInstance(w http.ResponseWriter, r *http.Request, ctx context.Context, project string) {
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

func (h *SpannerHandler) handleListDatabases(w http.ResponseWriter, r *http.Request, ctx context.Context, instanceName string) {
	items, _ := h.service.ListDatabases(ctx, instanceName)
	res := model.DatabasesList{
		Databases: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *SpannerHandler) handleCreateDatabase(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	var req struct {
		CreateStatement string `json:"createStatement"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	// Extract db name from statement
	dbName := parent + "/databases/test-db"
	db, _ := h.service.CreateDatabase(ctx, dbName)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(db)
}

func (h *SpannerHandler) handleCreateSession(w http.ResponseWriter, r *http.Request, ctx context.Context, dbName string) {
	session, _ := h.service.CreateSession(ctx, dbName)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func (h *SpannerHandler) handleExecuteSql(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	res := model.ResultSet{
		Metadata: &model.ResultSetMetadata{
			RowType: &model.StructType{
				Fields: []*model.Field{},
			},
		},
		Rows: [][]any{},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
