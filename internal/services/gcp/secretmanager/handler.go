package secretmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/secretmanager/model"
)

type SecretManagerGcpHandler struct {
	service *SecretManagerGcpService
}

func NewSecretManagerGcpHandler(service *SecretManagerGcpService) *SecretManagerGcpHandler {
	return &SecretManagerGcpHandler{
		service: service,
	}
}

func (h *SecretManagerGcpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/")
	ctx := context.Background()

	if strings.Contains(path, "/secrets") {
		if strings.HasSuffix(path, "/secrets") {
			if r.Method == "GET" {
				h.handleListSecrets(w, r, ctx)
				return
			}
			if r.Method == "POST" {
				h.handleCreateSecret(w, r, ctx)
				return
			}
		}

		if strings.Contains(path, ":addVersion") {
			secretName := strings.TrimSuffix(path, ":addVersion")
			h.handleAddVersion(w, r, ctx, secretName)
			return
		}

		if strings.Contains(path, ":access") {
			versionName := strings.TrimSuffix(path, ":access")
			h.handleAccessVersion(w, r, ctx, versionName)
			return
		}

		// Get/Delete
		if r.Method == "GET" {
			h.handleGetSecret(w, r, ctx, path)
			return
		}
		if r.Method == "DELETE" {
			h.handleDeleteSecret(w, r, ctx, path)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "Secret Manager path not implemented: %s"}}`, path)
}

func (h *SecretManagerGcpHandler) handleListSecrets(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	secrets, _ := h.service.ListSecrets(ctx)
	res := model.SecretsList{
		Secrets: secrets,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *SecretManagerGcpHandler) handleCreateSecret(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	var req struct {
		SecretId string `json:"secretId"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	
	// Path is like projects/P/secrets
	name := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/v1/"), "/") + "/" + req.SecretId
	
	secret, err := h.service.CreateSecret(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(secret)
}

func (h *SecretManagerGcpHandler) handleGetSecret(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	secret, err := h.service.GetSecret(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(secret)
}

func (h *SecretManagerGcpHandler) handleDeleteSecret(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	h.service.DeleteSecret(ctx, name)
	w.WriteHeader(http.StatusNoContent)
}

func (h *SecretManagerGcpHandler) handleAddVersion(w http.ResponseWriter, r *http.Request, ctx context.Context, secretName string) {
	var req struct {
		Payload *model.Payload `json:"payload"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	
	version, err := h.service.AddSecretVersion(ctx, secretName, req.Payload.Data)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version)
}

func (h *SecretManagerGcpHandler) handleAccessVersion(w http.ResponseWriter, r *http.Request, ctx context.Context, versionName string) {
	version, payload, err := h.service.AccessSecretVersion(ctx, versionName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	
	res := model.AccessSecretVersionResponse{
		Name: version.Name,
		Payload: &model.Payload{
			Data: payload,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
