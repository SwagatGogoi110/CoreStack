package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/iam/model"
)

type IamGcpHandler struct {
	service *IamGcpService
}

func NewIamGcpHandler(service *IamGcpService) *IamGcpHandler {
	return &IamGcpHandler{
		service: service,
	}
}

func (h *IamGcpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/")
	ctx := context.Background()

	if strings.Contains(path, "/serviceAccounts") {
		if strings.HasSuffix(path, "/serviceAccounts") {
			if r.Method == "GET" {
				project := strings.Split(path, "/")[1]
				h.handleListServiceAccounts(w, r, ctx, project)
				return
			}
			if r.Method == "POST" {
				project := strings.Split(path, "/")[1]
				h.handleCreateServiceAccount(w, r, ctx, project)
				return
			}
		}

		// Delete
		if r.Method == "DELETE" {
			h.handleDeleteServiceAccount(w, r, ctx, path)
			return
		}
		// Get
		if r.Method == "GET" {
			h.handleGetServiceAccount(w, r, ctx, path)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "GCP IAM path not implemented: %s"}}`, path)
}

func (h *IamGcpHandler) handleListServiceAccounts(w http.ResponseWriter, r *http.Request, ctx context.Context, project string) {
	accounts, _ := h.service.ListServiceAccounts(ctx, project)
	res := model.ServiceAccountsList{
		Accounts: accounts,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *IamGcpHandler) handleCreateServiceAccount(w http.ResponseWriter, r *http.Request, ctx context.Context, project string) {
	var req struct {
		AccountId      string `json:"accountId"`
		ServiceAccount struct {
			DisplayName string `json:"displayName"`
		} `json:"serviceAccount"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	
	sa, err := h.service.CreateServiceAccount(ctx, project, req.AccountId, req.ServiceAccount.DisplayName)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sa)
}

func (h *IamGcpHandler) handleGetServiceAccount(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	accounts, _ := h.service.ListServiceAccounts(ctx, "") // Scan all
	for _, sa := range accounts {
		if sa.Name == name {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(sa)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func (h *IamGcpHandler) handleDeleteServiceAccount(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	h.service.DeleteServiceAccount(ctx, name)
	w.WriteHeader(http.StatusNoContent)
}
