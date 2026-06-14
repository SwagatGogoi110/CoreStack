package artifactregistry

import (
	"context"
	"fmt"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/gcp/artifactregistry/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type ArtifactRegistryService struct {
	repoStore   storage.Backend[string, *model.Repository]
	imageStore  storage.Backend[string, *model.DockerImage]
}

func NewArtifactRegistryService(factory *storage.Factory) (*ArtifactRegistryService, error) {
	repoStore, _ := storage.CreateAccountAware[*model.Repository](factory, "artifactregistry", "ar-repos.json", "wal")
	imageStore, _ := storage.CreateAccountAware[*model.DockerImage](factory, "artifactregistry", "ar-images.json", "wal")

	return &ArtifactRegistryService{
		repoStore:   repoStore,
		imageStore:  imageStore,
	}, nil
}

// Repositories

func (s *ArtifactRegistryService) CreateRepository(ctx context.Context, name string, repo *model.Repository) (*model.Repository, error) {
	if _, ok, _ := s.repoStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("AlreadyExists: Repository %s already exists", name)
	}

	repo.Name = name
	repo.CreateTime = time.Now()
	repo.UpdateTime = time.Now()

	if err := s.repoStore.Put(ctx, name, repo); err != nil {
		return nil, err
	}
	return repo, nil
}

func (s *ArtifactRegistryService) ListRepositories(ctx context.Context, parent string) ([]*model.Repository, error) {
	return s.repoStore.Scan(ctx, func(k string) bool { return true })
}

// Docker Images

func (s *ArtifactRegistryService) ListDockerImages(ctx context.Context, parent string) ([]*model.DockerImage, error) {
	return s.imageStore.Scan(ctx, func(k string) bool { return true })
}
