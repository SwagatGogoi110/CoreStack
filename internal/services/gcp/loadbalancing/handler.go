package loadbalancing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/loadbalancing/model"
)

type LoadBalancingHandler struct {
	service *LoadBalancingService
}

func NewLoadBalancingHandler(service *LoadBalancingService) *LoadBalancingHandler {
	return &LoadBalancingHandler{
		service: service,
	}
}

func (h *LoadBalancingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	ctx := context.Background()

	if strings.Contains(path, "/forwardingRules") {
		if r.Method == "GET" {
			h.handleListForwardingRules(w, r, ctx)
			return
		}
		if r.Method == "POST" {
			h.handleCreateForwardingRule(w, r, ctx)
			return
		}
	}
	if strings.Contains(path, "/targetHttpProxies") {
		if r.Method == "GET" {
			h.handleListTargetHttpProxies(w, r, ctx)
			return
		}
		if r.Method == "POST" {
			h.handleCreateTargetHttpProxy(w, r, ctx)
			return
		}
	}
	if strings.Contains(path, "/urlMaps") {
		if r.Method == "GET" {
			h.handleListUrlMaps(w, r, ctx)
			return
		}
		if r.Method == "POST" {
			h.handleCreateUrlMap(w, r, ctx)
			return
		}
	}
	if strings.Contains(path, "/backendServices") {
		if r.Method == "GET" {
			h.handleListBackendServices(w, r, ctx)
			return
		}
		if r.Method == "POST" {
			h.handleCreateBackendService(w, r, ctx)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "Load Balancing path not implemented: %s"}}`, path)
}

func (h *LoadBalancingHandler) handleListForwardingRules(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	items, _ := h.service.ListForwardingRules(ctx)
	res := model.ForwardingRulesList{Items: items}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *LoadBalancingHandler) handleCreateForwardingRule(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	var fr model.ForwardingRule
	json.NewDecoder(r.Body).Decode(&fr)
	created, _ := h.service.CreateForwardingRule(ctx, fr.Name, &fr)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *LoadBalancingHandler) handleListTargetHttpProxies(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	items, _ := h.service.ListTargetHttpProxies(ctx)
	res := model.TargetHttpProxiesList{Items: items}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *LoadBalancingHandler) handleCreateTargetHttpProxy(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	var proxy model.TargetHttpProxy
	json.NewDecoder(r.Body).Decode(&proxy)
	created, _ := h.service.CreateTargetHttpProxy(ctx, proxy.Name, &proxy)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *LoadBalancingHandler) handleListUrlMaps(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	items, _ := h.service.ListUrlMaps(ctx)
	res := model.UrlMapsList{Items: items}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *LoadBalancingHandler) handleCreateUrlMap(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	var um model.UrlMap
	json.NewDecoder(r.Body).Decode(&um)
	created, _ := h.service.CreateUrlMap(ctx, um.Name, &um)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *LoadBalancingHandler) handleListBackendServices(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	items, _ := h.service.ListBackendServices(ctx)
	res := model.BackendServicesList{Items: items}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *LoadBalancingHandler) handleCreateBackendService(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	var be model.BackendService
	json.NewDecoder(r.Body).Decode(&be)
	created, _ := h.service.CreateBackendService(ctx, be.Name, &be)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}
