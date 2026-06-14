package spanner

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/gcp/spanner/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type SpannerService struct {
	instanceStore storage.Backend[string, *model.Instance]
	dbStore       storage.Backend[string, *model.Database]
	sessionStore  storage.Backend[string, *model.Session]
}

func NewSpannerService(factory *storage.Factory) (*SpannerService, error) {
	instanceStore, _ := storage.CreateAccountAware[*model.Instance](factory, "spanner", "spanner-instances.json", "wal")
	dbStore, _ := storage.CreateAccountAware[*model.Database](factory, "spanner", "spanner-databases.json", "wal")
	sessionStore, _ := storage.CreateAccountAware[*model.Session](factory, "spanner", "spanner-sessions.json", "wal")

	return &SpannerService{
		instanceStore: instanceStore,
		dbStore:       dbStore,
		sessionStore:  sessionStore,
	}, nil
}

// Instances

func (s *SpannerService) CreateInstance(ctx context.Context, name string, inst *model.Instance) (*model.Instance, error) {
	if _, ok, _ := s.instanceStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("AlreadyExists: Instance %s already exists", name)
	}

	inst.Name = name
	inst.State = "READY"
	inst.CreateTime = time.Now()
	inst.UpdateTime = time.Now()

	if err := s.instanceStore.Put(ctx, name, inst); err != nil {
		return nil, err
	}
	return inst, nil
}

func (s *SpannerService) ListInstances(ctx context.Context, project string) ([]*model.Instance, error) {
	return s.instanceStore.Scan(ctx, func(k string) bool { return true })
}

// Databases

func (s *SpannerService) CreateDatabase(ctx context.Context, name string) (*model.Database, error) {
	db := &model.Database{
		Name:       name,
		State:      "READY",
		CreateTime: time.Now(),
	}
	if err := s.dbStore.Put(ctx, name, db); err != nil {
		return nil, err
	}
	return db, nil
}

func (s *SpannerService) ListDatabases(ctx context.Context, instanceName string) ([]*model.Database, error) {
	return s.dbStore.Scan(ctx, func(k string) bool { return true })
}

// Sessions

func (s *SpannerService) CreateSession(ctx context.Context, dbName string) (*model.Session, error) {
	id := uuid.New().String()
	name := dbName + "/sessions/" + id
	session := &model.Session{
		Name:                   name,
		CreateTime:             time.Now(),
		ApproximateLastUseTime: time.Now(),
	}
	s.sessionStore.Put(ctx, name, session)
	return session, nil
}
func (s *SpannerService) DeleteSession(ctx context.Context, name string) error {
	return s.sessionStore.Delete(ctx, name)
}
