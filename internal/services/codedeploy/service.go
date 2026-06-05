package codedeploy

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/codedeploy/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CodeDeployService struct {
	appStore   storage.Backend[string, *model.Application]
	groupStore storage.Backend[string, *model.DeploymentGroup]
	mu         sync.RWMutex
}

func NewCodeDeployService(factory *storage.Factory) (*CodeDeployService, error) {
	appStore, _ := storage.CreateAccountAware[*model.Application](factory, "codedeploy", "cd-apps.json", "wal")
	groupStore, _ := storage.CreateAccountAware[*model.DeploymentGroup](factory, "codedeploy", "cd-groups.json", "wal")

	return &CodeDeployService{
		appStore:   appStore,
		groupStore: groupStore,
	}, nil
}

func (s *CodeDeployService) CreateApplication(ctx context.Context, name, platform string) (*model.Application, error) {
	if _, ok, _ := s.appStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("ApplicationAlreadyExistsException: Application already exists")
	}

	app := &model.Application{
		ApplicationID:   uuid.New().String(),
		ApplicationName: name,
		CreateTime:      time.Now(),
		ComputePlatform: platform,
	}

	if err := s.appStore.Put(ctx, name, app); err != nil {
		return nil, err
	}

	log.Printf("Created CodeDeploy application: %s", name)
	return app, nil
}

func (s *CodeDeployService) GetApplication(ctx context.Context, name string) (*model.Application, error) {
	app, ok, err := s.appStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("ApplicationDoesNotExistException: Application not found")
	}
	return app, nil
}

func (s *CodeDeployService) CreateDeploymentGroup(ctx context.Context, appName, groupName string) (*model.DeploymentGroup, error) {
	s.GetApplication(ctx, appName)
	
	key := fmt.Sprintf("%s:%s", appName, groupName)
	if _, ok, _ := s.groupStore.Get(ctx, key); ok {
		return nil, fmt.Errorf("DeploymentGroupAlreadyExistsException: Deployment group already exists")
	}

	group := &model.DeploymentGroup{
		ApplicationName:     appName,
		DeploymentGroupID:   uuid.New().String(),
		DeploymentGroupName: groupName,
	}

	if err := s.groupStore.Put(ctx, key, group); err != nil {
		return nil, err
	}

	log.Printf("Created CodeDeploy group: %s for app %s", groupName, appName)
	return group, nil
}

func (s *CodeDeployService) ListApplications(ctx context.Context) ([]string, error) {
	apps, err := s.appStore.Scan(ctx, func(k string) bool { return true })
	if err != nil {
		return nil, err
	}
	var names []string
	for _, a := range apps {
		names = append(names, a.ApplicationName)
	}
	return names, nil
}
