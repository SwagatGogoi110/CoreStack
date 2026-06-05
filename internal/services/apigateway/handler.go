package apigateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type ApiGatewayHandler struct {
	service *ApiGatewayService
}

func NewApiGatewayHandler(service *ApiGatewayService) *ApiGatewayHandler {
	return &ApiGatewayHandler{
		service: service,
	}
}

func (h *ApiGatewayHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("ApiGateway does not support standard JSON protocol dispatcher")
}

func (h *ApiGatewayHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("ApiGateway does not support Query protocol")
}

func (h *ApiGatewayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/restapis")
	
	if path == "" || path == "/" {
		if r.Method == "POST" {
			h.handleCreateRestApi(w, r)
			return
		}
		if r.Method == "GET" {
			h.handleListRestApis(w, r)
			return
		}
	}

	if strings.HasPrefix(path, "/") {
		parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
		apiID := parts[0]
		
		if len(parts) == 1 && r.Method == "GET" {
			h.handleGetRestApi(w, r, apiID)
			return
		}
		
		if len(parts) >= 2 && parts[1] == "resources" {
			if len(parts) == 2 && r.Method == "GET" {
				h.handleGetResources(w, r, apiID)
				return
			}
			// More complex resource routing...
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func (h *ApiGatewayHandler) handleCreateRestApi(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	ctx := context.Background()
	api, err := h.service.CreateRestApi(ctx, body.Name, body.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(api)
}

func (h *ApiGatewayHandler) handleListRestApis(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *ApiGatewayHandler) handleGetRestApi(w http.ResponseWriter, r *http.Request, apiID string) {
	ctx := context.Background()
	api, err := h.service.GetRestApi(ctx, apiID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(api)
}

func (h *ApiGatewayHandler) handleGetResources(w http.ResponseWriter, r *http.Request, apiID string) {
	ctx := context.Background()
	resources, _ := h.service.GetResources(ctx, apiID)
	res := map[string]any{
		"items": resources,
	}
	json.NewEncoder(w).Encode(res)
}
