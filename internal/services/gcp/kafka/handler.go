package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/kafka/model"
)

type KafkaHandler struct {
	service *KafkaService
}

func NewKafkaHandler(service *KafkaService) *KafkaHandler {
	return &KafkaHandler{
		service: service,
	}
}

func (h *KafkaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

		if r.Method == "GET" {
			h.handleGetCluster(w, r, ctx, path)
			return
		}
		if r.Method == "DELETE" {
			h.handleDeleteCluster(w, r, ctx, path)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "GCP Kafka path not implemented: %s"}}`, path)
}

func (h *KafkaHandler) handleListClusters(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	clusters, _ := h.service.ListClusters(ctx, parent)
	res := model.ClustersList{
		Clusters: clusters,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *KafkaHandler) handleCreateCluster(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	var c model.Cluster
	json.NewDecoder(r.Body).Decode(&c)
	
	clusterId := r.URL.Query().Get("clusterId")
	name := parent + "/" + clusterId

	created, _ := h.service.CreateCluster(ctx, name, &c)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *KafkaHandler) handleGetCluster(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	c, err := h.service.GetCluster(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func (h *KafkaHandler) handleDeleteCluster(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	h.service.DeleteCluster(ctx, name)
	w.WriteHeader(http.StatusNoContent)
}
