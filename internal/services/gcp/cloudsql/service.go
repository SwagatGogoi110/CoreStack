package cloudsql

import (
	"context"
	"fmt"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/gcp/cloudsql/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CloudSqlService struct {
	instanceStore storage.Backend[string, *model.DatabaseInstance]
	dbStore       storage.Backend[string, *model.Database]
	userStore     storage.Backend[string, *model.User]
}

func NewCloudSqlService(factory *storage.Factory) (*CloudSqlService, error) {
	instanceStore, _ := storage.CreateAccountAware[*model.DatabaseInstance](factory, "cloudsql", "sql-instances.json", "wal")
	dbStore, _ := storage.CreateAccountAware[*model.Database](factory, "cloudsql", "sql-databases.json", "wal")
	userStore, _ := storage.CreateAccountAware[*model.User](factory, "cloudsql", "sql-users.json", "wal")

	return &CloudSqlService{
		instanceStore: instanceStore,
		dbStore:       dbStore,
		userStore:     userStore,
	}, nil
}

// Instances

func (s *CloudSqlService) CreateInstance(ctx context.Context, project string, inst *model.DatabaseInstance) (*model.DatabaseInstance, error) {
	key := fmt.Sprintf("%s:%s", project, inst.Name)
	if _, ok, _ := s.instanceStore.Get(ctx, key); ok {
		return nil, fmt.Errorf("AlreadyExists: Instance %s already exists", inst.Name)
	}

	inst.Project = project
	inst.Kind = "sql#instance"
	inst.State = "RUNNABLE"
	inst.CreateTime = time.Now()
	inst.IpAddresses = []*model.IpMapping{
		{Type: "PRIMARY", IpAddress: "127.0.0.1"},
	}

	if err := s.instanceStore.Put(ctx, key, inst); err != nil {
		return nil, err
	}

	return inst, nil
}

func (s *CloudSqlService) GetInstance(ctx context.Context, project, name string) (*model.DatabaseInstance, error) {
	key := fmt.Sprintf("%s:%s", project, name)
	inst, ok, err := s.instanceStore.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NotFound: Instance %s not found", name)
	}
	return inst, nil
}

func (s *CloudSqlService) ListInstances(ctx context.Context, project string) ([]*model.DatabaseInstance, error) {
	return s.instanceStore.Scan(ctx, func(k string) bool {
		inst, _, _ := s.instanceStore.Get(ctx, k)
		return inst.Project == project
	})
}

func (s *CloudSqlService) DeleteInstance(ctx context.Context, project, name string) error {
	key := fmt.Sprintf("%s:%s", project, name)
	return s.instanceStore.Delete(ctx, key)
}

// Databases

func (s *CloudSqlService) CreateDatabase(ctx context.Context, project, instance, dbName string) (*model.Database, error) {
	key := fmt.Sprintf("%s:%s:%s", project, instance, dbName)
	db := &model.Database{
		Kind:     "sql#database",
		Name:     dbName,
		Project:  project,
		Instance: instance,
	}
	if err := s.dbStore.Put(ctx, key, db); err != nil {
		return nil, err
	}
	return db, nil
}

func (s *CloudSqlService) ListDatabases(ctx context.Context, project, instance string) ([]*model.Database, error) {
	return s.dbStore.Scan(ctx, func(k string) bool {
		db, _, _ := s.dbStore.Get(ctx, k)
		return db.Project == project && db.Instance == instance
	})
}

// Users

func (s *CloudSqlService) CreateUser(ctx context.Context, project, instance, userName string) (*model.User, error) {
	key := fmt.Sprintf("%s:%s:%s", project, instance, userName)
	user := &model.User{
		Kind:     "sql#user",
		Name:     userName,
		Project:  project,
		Instance: instance,
	}
	if err := s.userStore.Put(ctx, key, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *CloudSqlService) ListUsers(ctx context.Context, project, instance string) ([]*model.User, error) {
	return s.userStore.Scan(ctx, func(k string) bool {
		user, _, _ := s.userStore.Get(ctx, k)
		return user.Project == project && user.Instance == instance
	})
}
