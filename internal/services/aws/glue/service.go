package glue

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/aws/glue/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type GlueService struct {
	dbStore    storage.Backend[string, *model.Database]
	tableStore storage.Backend[string, *model.Table]
	mu         sync.RWMutex
}

func NewGlueService(factory *storage.Factory) (*GlueService, error) {
	dbStore, _ := storage.CreateAccountAware[*model.Database](factory, "glue", "glue-databases.json", "wal")
	tableStore, _ := storage.CreateAccountAware[*model.Table](factory, "glue", "glue-tables.json", "wal")

	return &GlueService{
		dbStore:    dbStore,
		tableStore: tableStore,
	}, nil
}

func (s *GlueService) CreateDatabase(ctx context.Context, name, description string) (*model.Database, error) {
	if _, ok, _ := s.dbStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("AlreadyExistsException: Database already exists")
	}

	db := &model.Database{
		Name:        name,
		Description: description,
		CreateTime:  time.Now(),
	}

	if err := s.dbStore.Put(ctx, name, db); err != nil {
		return nil, err
	}

	log.Printf("Created Glue database: %s", name)
	return db, nil
}

func (s *GlueService) GetDatabase(ctx context.Context, name string) (*model.Database, error) {
	db, ok, err := s.dbStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("EntityNotFoundException: Database not found")
	}
	return db, nil
}

func (s *GlueService) CreateTable(ctx context.Context, dbName string, table *model.Table) error {
	s.GetDatabase(ctx, dbName)

	key := fmt.Sprintf("%s:%s", dbName, table.Name)
	if _, ok, _ := s.tableStore.Get(ctx, key); ok {
		return fmt.Errorf("AlreadyExistsException: Table already exists")
	}

	table.DatabaseName = dbName
	table.CreateTime = time.Now()
	table.UpdateTime = time.Now()

	if err := s.tableStore.Put(ctx, key, table); err != nil {
		return err
	}

	log.Printf("Created Glue table: %s.%s", dbName, table.Name)
	return nil
}

func (s *GlueService) GetTables(ctx context.Context, dbName string) ([]*model.Table, error) {
	prefix := dbName + ":"
	return s.tableStore.Scan(ctx, func(k string) bool {
		return len(k) > len(prefix) && k[:len(prefix)] == prefix
	})
}
