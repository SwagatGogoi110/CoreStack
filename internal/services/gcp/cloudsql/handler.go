package cloudsql

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/cloudsql/model"
)

type CloudSqlHandler struct {
	service *CloudSqlService
}

func NewCloudSqlHandler(service *CloudSqlService) *CloudSqlHandler {
	return &CloudSqlHandler{
		service: service,
	}
}

func (h *CloudSqlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/sql/v1/")
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

		instance := parts[3]
		if len(parts) == 4 {
			if r.Method == "GET" {
				h.handleGetInstance(w, r, ctx, project, instance)
				return
			}
			if r.Method == "DELETE" {
				h.handleDeleteInstance(w, r, ctx, project, instance)
				return
			}
		}

		if len(parts) == 5 {
			if parts[4] == "databases" {
				if r.Method == "GET" {
					h.handleListDatabases(w, r, ctx, project, instance)
					return
				}
				if r.Method == "POST" {
					h.handleCreateDatabase(w, r, ctx, project, instance)
					return
				}
			}
			if parts[4] == "users" {
				if r.Method == "GET" {
					h.handleListUsers(w, r, ctx, project, instance)
					return
				}
				if r.Method == "POST" {
					h.handleCreateUser(w, r, ctx, project, instance)
					return
				}
			}
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "Cloud SQL path not implemented: %s"}}`, path)
}

func (h *CloudSqlHandler) handleListInstances(w http.ResponseWriter, r *http.Request, ctx context.Context, project string) {
	items, _ := h.service.ListInstances(ctx, project)
	res := model.InstancesList{
		Kind:  "sql#instancesList",
		Items: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *CloudSqlHandler) handleCreateInstance(w http.ResponseWriter, r *http.Request, ctx context.Context, project string) {
	var req model.DatabaseInstance
	json.NewDecoder(r.Body).Decode(&req)
	
	created, err := h.service.CreateInstance(ctx, project, &req)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *CloudSqlHandler) handleGetInstance(w http.ResponseWriter, r *http.Request, ctx context.Context, project, name string) {
	inst, err := h.service.GetInstance(ctx, project, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inst)
}

func (h *CloudSqlHandler) handleDeleteInstance(w http.ResponseWriter, r *http.Request, ctx context.Context, project, name string) {
	h.service.DeleteInstance(ctx, project, name)
	w.WriteHeader(http.StatusNoContent)
}

func (h *CloudSqlHandler) handleListDatabases(w http.ResponseWriter, r *http.Request, ctx context.Context, project, instance string) {
	items, _ := h.service.ListDatabases(ctx, project, instance)
	res := model.DatabasesList{
		Kind:  "sql#databasesList",
		Items: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *CloudSqlHandler) handleCreateDatabase(w http.ResponseWriter, r *http.Request, ctx context.Context, project, instance string) {
	var req struct {
		Name string `json:"name"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	db, _ := h.service.CreateDatabase(ctx, project, instance, req.Name)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(db)
}

func (h *CloudSqlHandler) handleListUsers(w http.ResponseWriter, r *http.Request, ctx context.Context, project, instance string) {
	items, _ := h.service.ListUsers(ctx, project, instance)
	res := model.UsersList{
		Kind:  "sql#usersList",
		Items: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *CloudSqlHandler) handleCreateUser(w http.ResponseWriter, r *http.Request, ctx context.Context, project, instance string) {
	var req struct {
		Name string `json:"name"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	user, _ := h.service.CreateUser(ctx, project, instance, req.Name)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
