package ecr

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/ecr/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type EcrService struct {
	repoStore storage.Backend[string, *model.Repository]
}

func NewEcrService(factory *storage.Factory) (*EcrService, error) {
	repoStore, _ := storage.CreateAccountAware[*model.Repository](factory, "ecr", "ecr-repositories.json", "wal")

	return &EcrService{
		repoStore: repoStore,
	}, nil
}

func (s *EcrService) CreateRepository(ctx context.Context, name string) (*model.Repository, error) {
	if _, ok, _ := s.repoStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("RepositoryAlreadyExistsException: Repository already exists")
	}

	arn := fmt.Sprintf("arn:aws:ecr:us-east-1:000000000000:repository/%s", name)
	repo := &model.Repository{
		RepositoryArn:      arn,
		RegistryID:         "000000000000",
		RepositoryName:     name,
		RepositoryUri:      fmt.Sprintf("localhost:8080/%s", name),
		CreatedAt:          time.Now(),
		ImageTagMutability: "MUTABLE",
	}

	if err := s.repoStore.Put(ctx, name, repo); err != nil {
		return nil, err
	}

	log.Printf("Created ECR repository: %s", name)
	return repo, nil
}

func (s *EcrService) DescribeRepositories(ctx context.Context, names []string) ([]*model.Repository, error) {
	return s.repoStore.Scan(ctx, func(k string) bool {
		if len(names) == 0 {
			return true
		}
		for _, n := range names {
			if n == k {
				return true
			}
		}
		return false
	})
}
