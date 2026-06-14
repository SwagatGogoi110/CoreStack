package container

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/container/model"
)

type GkeHandler struct {
	service *GkeService
}

func NewGkeHandler(service *GkeService) *GkeHandler {
	return &GkeHandler{
		service: service,
	}
}

func (h *GkeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/")
	ctx := context.Background()

	if strings.Contains(path, "/clusters") {
		
		if strings.HasSuffix(path, "/clusters") {
			if r.Method == "GET" {
				h.handleListClusters(w, r, ctx, path)
				return
			}
			if r.Method == "POST" {
				h.handleCreateCluster(w, r, ctx, path)
				return
			}
		}

		if strings.Contains(path, "/nodePools") {
			if strings.HasSuffix(path, "/nodePools") {
				if r.Method == "GET" {
					h.handleListNodePools(w, r, ctx, path)
					return
				}
				if r.Method == "POST" {
					h.handleCreateNodePool(w, r, ctx, path)
					return
				}
			}
		}

		if r.Method == "GET" {
			h.handleGetCluster(w, r, ctx, path)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "GKE path not implemented: %s"}}`, path)
}

func (h *GkeHandler) handleListClusters(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	items, _ := h.service.ListClusters(ctx, parent)
	res := model.ClustersList{
		Clusters: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *GkeHandler) handleCreateCluster(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	var req struct {
		Cluster *model.Cluster `json:"cluster"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	name := parent + "/" + req.Cluster.Name
	created, _ := h.service.CreateCluster(ctx, name, req.Cluster)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *GkeHandler) handleGetCluster(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	clusters, _ := h.service.ListClusters(ctx, "")
	for _, c := range clusters {
		if c.Name == name {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func (h *GkeHandler) handleListNodePools(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	items, _ := h.service.ListNodePools(ctx, parent)
	res := model.NodePoolsList{
		NodePools: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *GkeHandler) handleCreateNodePool(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	var req struct {
		NodePool *model.NodePool `json:"nodePool"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	created, _ := h.service.CreateNodePool(ctx, parent, req.NodePool.Name, req.NodePool)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}
