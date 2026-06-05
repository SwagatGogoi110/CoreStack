package codebuild

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/codebuild/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CodeBuildService struct {
	projectStore storage.Backend[string, *model.Project]
	buildStore   storage.Backend[string, *model.Build]
	mu           sync.RWMutex
}

func NewCodeBuildService(factory *storage.Factory) (*CodeBuildService, error) {
	projectStore, _ := storage.CreateAccountAware[*model.Project](factory, "codebuild", "cb-projects.json", "wal")
	buildStore, _ := storage.CreateAccountAware[*model.Build](factory, "codebuild", "cb-builds.json", "wal")

	return &CodeBuildService{
		projectStore: projectStore,
		buildStore:   buildStore,
	}, nil
}

func (s *CodeBuildService) CreateProject(ctx context.Context, name, description string) (*model.Project, error) {
	if _, ok, _ := s.projectStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("ResourceAlreadyExistsException: Project already exists")
	}

	arn := fmt.Sprintf("arn:aws:codebuild:us-east-1:000000000000:project/%s", name)
	project := &model.Project{
		Name:         name,
		Arn:          arn,
		Description:  description,
		Created:      time.Now(),
		LastModified: time.Now(),
	}

	if err := s.projectStore.Put(ctx, name, project); err != nil {
		return nil, err
	}

	log.Printf("Created CodeBuild project: %s", name)
	return project, nil
}

func (s *CodeBuildService) BatchGetProjects(ctx context.Context, names []string) ([]*model.Project, error) {
	var projects []*model.Project
	for _, n := range names {
		if p, ok, _ := s.projectStore.Get(ctx, n); ok {
			projects = append(projects, p)
		}
	}
	return projects, nil
}

func (s *CodeBuildService) ListProjects(ctx context.Context) ([]string, error) {
	projects, err := s.projectStore.Scan(ctx, func(k string) bool { return true })
	if err != nil {
		return nil, err
	}
	var names []string
	for _, p := range projects {
		names = append(names, p.Name)
	}
	return names, nil
}
