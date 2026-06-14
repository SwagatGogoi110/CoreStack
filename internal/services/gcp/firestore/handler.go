package firestore

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/firestore/model"
)

type FirestoreHandler struct {
	service *FirestoreService
}

func NewFirestoreHandler(service *FirestoreService) *FirestoreHandler {
	return &FirestoreHandler{
		service: service,
	}
}

func (h *FirestoreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/")
	ctx := context.Background()

	// Pattern: projects/{project}/databases/{database}/documents/{collection}/{document}
	if strings.Contains(path, "/documents") {
		parts := strings.Split(path, "/")
		
		// Find index of "documents"
		idx := -1
		for i, p := range parts {
			if p == "documents" {
				idx = i
				break
			}
		}

		if idx != -1 {
			if r.Method == "GET" {
				if idx == len(parts)-2 {
					// List collection: .../documents/{collection}
					h.handleListDocuments(w, r, ctx, path)
					return
				}
				if idx < len(parts)-2 {
					// Get document: .../documents/{collection}/{id}
					h.handleGetDocument(w, r, ctx, path)
					return
				}
			}
			if r.Method == "POST" && idx == len(parts)-2 {
				// Create document: POST .../documents/{collection}?documentId=ID
				h.handleCreateDocument(w, r, ctx, path)
				return
			}
			if r.Method == "PATCH" {
				// Update document: PATCH .../documents/{collection}/{id}
				h.handleCreateDocument(w, r, ctx, path) // Reuse create for stub
				return
			}
			if r.Method == "DELETE" {
				h.handleDeleteDocument(w, r, ctx, path)
				return
			}
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "Firestore path not implemented: %s"}}`, path)
}

func (h *FirestoreHandler) handleListDocuments(w http.ResponseWriter, r *http.Request, ctx context.Context, collectionPath string) {
	docs, _ := h.service.ListDocuments(ctx, collectionPath)
	res := model.DocumentsList{
		Documents: docs,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *FirestoreHandler) handleGetDocument(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	doc, err := h.service.GetDocument(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

func (h *FirestoreHandler) handleCreateDocument(w http.ResponseWriter, r *http.Request, ctx context.Context, path string) {
	var doc model.Document
	json.NewDecoder(r.Body).Decode(&doc)
	
	name := doc.Name
	if name == "" {
		docId := r.URL.Query().Get("documentId")
		if docId != "" {
			name = path + "/" + docId
		} else {
			name = path + "/dummy-id"
		}
	}

	created, err := h.service.CreateDocument(ctx, name, doc.Fields)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *FirestoreHandler) handleDeleteDocument(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	h.service.DeleteDocument(ctx, name)
	w.WriteHeader(http.StatusNoContent)
}
