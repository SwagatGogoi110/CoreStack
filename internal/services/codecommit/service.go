package codecommit

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/codecommit/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CodecommitService struct {
	repoStore storage.Backend[string, *model.Repository]
}

func NewCodecommitService(factory *storage.Factory) (*CodecommitService, error) {
	repoStore, _ := storage.CreateAccountAware[*model.Repository](factory, "codecommit", "codecommit-repos.json", "wal")

	return &CodecommitService{
		repoStore: repoStore,
	}, nil
}

func (s *CodecommitService) CreateRepository(ctx context.Context, name, description string) (*model.Repository, error) {
	if _, ok, _ := s.repoStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("RepositoryNameExistsException: Repository name already exists")
	}

	repo := &model.Repository{
		AccountId:             "000000000000",
		RepositoryId:          uuid.New().String(),
		RepositoryName:        name,
		RepositoryDescription: description,
		CreationDate:          time.Now(),
		LastModifiedDate:      time.Now(),
		Arn:                   fmt.Sprintf("arn:aws:codecommit:us-east-1:000000000000:%s", name),
		CloneUrlHttp:          fmt.Sprintf("http://localhost:4566/codecommit/%s", name),
	}

	if err := s.repoStore.Put(ctx, name, repo); err != nil {
		return nil, err
	}

	return repo, nil
}

func (s *CodecommitService) GetRepository(ctx context.Context, name string) (*model.Repository, error) {
	repo, ok, err := s.repoStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("RepositoryDoesNotExistException: Repository %s does not exist", name)
	}
	return repo, nil
}

func (s *CodecommitService) ListRepositories(ctx context.Context) ([]*model.Repository, error) {
	return s.repoStore.Scan(ctx, func(k string) bool { return true })
}

func (s *CodecommitService) DeleteRepository(ctx context.Context, name string) error {
	return s.repoStore.Delete(ctx, name)
}
