package datastore

import (
	"context"

	"github.com/hectorvent/cloudstack/internal/services/gcp/datastore/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type DatastoreService struct {
	entityStore storage.Backend[string, *model.Entity]
}

func NewDatastoreService(factory *storage.Factory) (*DatastoreService, error) {
	entityStore, _ := storage.CreateAccountAware[*model.Entity](factory, "datastore", "ds-entities.json", "wal")

	return &DatastoreService{
		entityStore: entityStore,
	}, nil
}

func (s *DatastoreService) Commit(ctx context.Context, project string, mutations []map[string]any) ([]*model.MutationResult, error) {
	results := make([]*model.MutationResult, 0)
	
	// Very simplified mutation logic for stub
	for _, m := range mutations {
		if _, ok := m["insert"].(map[string]any); ok {
			// Port logic from Java to handle Entity conversion if needed
			// For stub, just return success
			results = append(results, &model.MutationResult{Version: "1"})
		}
	}
	
	return results, nil
}

func (s *DatastoreService) Lookup(ctx context.Context, keys []*model.Key) (*model.LookupResponse, error) {
	res := &model.LookupResponse{
		Found:   []*model.EntityResult{},
		Missing: []*model.EntityResult{},
	}
	
	for _, k := range keys {
		// keyToString helper
		keyStr := s.keyToString(k)
		entity, ok, _ := s.entityStore.Get(ctx, keyStr)
		if ok {
			res.Found = append(res.Found, &model.EntityResult{Entity: entity})
		} else {
			res.Missing = append(res.Missing, &model.EntityResult{Entity: &model.Entity{Key: k}})
		}
	}
	
	return res, nil
}

func (s *DatastoreService) keyToString(k *model.Key) string {
	// Simplified key serialization
	res := ""
	if k.PartitionId != nil {
		res += k.PartitionId.NamespaceId + ":"
	}
	for _, p := range k.Path {
		res += p.Kind + ":" + p.Name + p.Id + "/"
	}
	return res
}
