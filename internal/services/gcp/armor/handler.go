package armor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/armor/model"
)

type CloudArmorHandler struct {
	service *CloudArmorService
}

func NewCloudArmorHandler(service *CloudArmorService) *CloudArmorHandler {
	return &CloudArmorHandler{
		service: service,
	}
}

func (h *CloudArmorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/compute/v1/projects/")
	ctx := context.Background()

	if strings.Contains(path, "/global/securityPolicies") {
		if strings.HasSuffix(path, "/global/securityPolicies") {
			if r.Method == "GET" {
				h.handleListPolicies(w, r, ctx)
				return
			}
			if r.Method == "POST" {
				h.handleCreatePolicy(w, r, ctx)
				return
			}
		}

		parts := strings.Split(path, "/")
		policyName := parts[len(parts)-1]
		if r.Method == "GET" {
			h.handleGetPolicy(w, r, ctx, policyName)
			return
		}
		if r.Method == "DELETE" {
			h.handleDeletePolicy(w, r, ctx, policyName)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "Cloud Armor path not implemented: %s"}}`, path)
}

func (h *CloudArmorHandler) handleListPolicies(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	items, _ := h.service.ListPolicies(ctx)
	res := model.SecurityPoliciesList{
		SecurityPolicies: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *CloudArmorHandler) handleCreatePolicy(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	var policy model.SecurityPolicy
	json.NewDecoder(r.Body).Decode(&policy)
	created, _ := h.service.CreatePolicy(ctx, policy.Name, &policy)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *CloudArmorHandler) handleGetPolicy(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	items, _ := h.service.ListPolicies(ctx)
	for _, p := range items {
		if p.Name == name {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func (h *CloudArmorHandler) handleDeletePolicy(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	h.service.DeletePolicy(ctx, name)
	w.WriteHeader(http.StatusNoContent)
}
