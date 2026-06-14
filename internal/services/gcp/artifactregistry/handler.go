package artifactregistry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/artifactregistry/model"
)

type ArtifactRegistryHandler struct {
	service *ArtifactRegistryService
}

func NewArtifactRegistryHandler(service *ArtifactRegistryService) *ArtifactRegistryHandler {
	return &ArtifactRegistryHandler{
		service: service,
	}
}

func (h *ArtifactRegistryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/")
	ctx := context.Background()

	if strings.Contains(path, "/repositories") {
		if strings.HasSuffix(path, "/repositories") {
			if r.Method == "GET" {
				h.handleListRepositories(w, r, ctx, path)
				return
			}
			if r.Method == "POST" {
				h.handleCreateRepository(w, r, ctx, path)
				return
			}
		}

		if strings.Contains(path, "/dockerImages") {
			if r.Method == "GET" {
				h.handleListDockerImages(w, r, ctx, path)
				return
			}
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "Artifact Registry path not implemented: %s"}}`, path)
}

func (h *ArtifactRegistryHandler) handleListRepositories(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	items, _ := h.service.ListRepositories(ctx, parent)
	res := model.RepositoriesList{
		Repositories: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *ArtifactRegistryHandler) handleCreateRepository(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	var repo model.Repository
	json.NewDecoder(r.Body).Decode(&repo)
	
	repositoryId := r.URL.Query().Get("repositoryId")
	name := parent + "/" + repositoryId

	created, _ := h.service.CreateRepository(ctx, name, &repo)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *ArtifactRegistryHandler) handleListDockerImages(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	items, _ := h.service.ListDockerImages(ctx, parent)
	res := model.DockerImagesList{
		DockerImages: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
