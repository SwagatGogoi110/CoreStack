package firestore

import (
	"context"
	"fmt"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/gcp/firestore/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type FirestoreService struct {
	docStore storage.Backend[string, *model.Document]
}

func NewFirestoreService(factory *storage.Factory) (*FirestoreService, error) {
	docStore, _ := storage.CreateAccountAware[*model.Document](factory, "firestore", "firestore-documents.json", "wal")

	return &FirestoreService{
		docStore: docStore,
	}, nil
}

func (s *FirestoreService) CreateDocument(ctx context.Context, name string, fields map[string]model.Value) (*model.Document, error) {
	doc := &model.Document{
		Name:       name,
		Fields:     fields,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := s.docStore.Put(ctx, name, doc); err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *FirestoreService) GetDocument(ctx context.Context, name string) (*model.Document, error) {
	doc, ok, err := s.docStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NotFound: Document %s not found", name)
	}
	return doc, nil
}

func (s *FirestoreService) ListDocuments(ctx context.Context, collectionId string) ([]*model.Document, error) {
	// collectionId is a path prefix like projects/P/databases/D/documents/C
	return s.docStore.Scan(ctx, func(k string) bool {
		// Basic prefix matching
		return len(k) > len(collectionId) && k[:len(collectionId)] == collectionId
	})
}

func (s *FirestoreService) DeleteDocument(ctx context.Context, name string) error {
	return s.docStore.Delete(ctx, name)
}
