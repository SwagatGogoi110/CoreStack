package gcs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/gcs/model"
)

type GcsHandler struct {
	service *GcsService
}

func NewGcsHandler(service *GcsService) *GcsHandler {
	return &GcsHandler{
		service: service,
	}
}

func (h *GcsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	ctx := context.Background()

	// GCS API Path patterns:
	// /storage/v1/b -> List Buckets
	// /storage/v1/b/BUCKET -> Get/Delete Bucket
	// /storage/v1/b/BUCKET/o -> List Objects
	// /storage/v1/b/BUCKET/o/OBJECT -> Get/Delete Object
	// /upload/storage/v1/b/BUCKET/o -> Upload Object

	if strings.HasPrefix(path, "/storage/v1/b") {
		subPath := strings.TrimPrefix(path, "/storage/v1/b")
		if subPath == "" || subPath == "/" {
			if r.Method == "GET" {
				h.handleListBuckets(w, r, ctx)
				return
			}
			if r.Method == "POST" {
				h.handleCreateBucket(w, r, ctx)
				return
			}
		}

		parts := strings.Split(strings.Trim(subPath, "/"), "/")
		bucketName := parts[0]

		if len(parts) == 1 {
			if r.Method == "GET" {
				h.handleGetBucket(w, r, ctx, bucketName)
				return
			}
			if r.Method == "DELETE" {
				h.handleDeleteBucket(w, r, ctx, bucketName)
				return
			}
		}

		if len(parts) >= 2 && parts[1] == "o" {
			if len(parts) == 2 {
				if r.Method == "GET" {
					h.handleListObjects(w, r, ctx, bucketName)
					return
				}
			} else {
				objectName := strings.Join(parts[2:], "/")
				if r.Method == "GET" {
					h.handleGetObject(w, r, ctx, bucketName, objectName)
					return
				}
				if r.Method == "DELETE" {
					h.handleDeleteObject(w, r, ctx, bucketName, objectName)
					return
				}
			}
		}
	}

	if strings.HasPrefix(path, "/upload/storage/v1/b/") {
		bucketName := strings.Split(strings.TrimPrefix(path, "/upload/storage/v1/b/"), "/")[0]
		if r.Method == "POST" {
			h.handleUploadObject(w, r, ctx, bucketName)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "GCS path not implemented: %s"}}`, path)
}

func (h *GcsHandler) handleListBuckets(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	buckets, _ := h.service.ListBuckets(ctx)
	res := model.BucketsList{
		Kind:  "storage#buckets",
		Items: buckets,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *GcsHandler) handleCreateBucket(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	var req struct {
		Name string `json:"name"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	bucket, err := h.service.CreateBucket(ctx, req.Name)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bucket)
}

func (h *GcsHandler) handleGetBucket(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	bucket, err := h.service.GetBucket(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bucket)
}

func (h *GcsHandler) handleDeleteBucket(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	err := h.service.DeleteBucket(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *GcsHandler) handleListObjects(w http.ResponseWriter, r *http.Request, ctx context.Context, bucketName string) {
	objects, _ := h.service.ListObjects(ctx, bucketName)
	res := model.ObjectsList{
		Kind:  "storage#objects",
		Items: objects,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *GcsHandler) handleGetObject(w http.ResponseWriter, r *http.Request, ctx context.Context, bucketName, objectName string) {
	alt := r.URL.Query().Get("alt")
	meta, content, err := h.service.GetObject(ctx, bucketName, objectName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer content.Close()

	if alt == "media" {
		w.Header().Set("Content-Type", meta.ContentType)
		w.Header().Set("Content-Length", meta.Size)
		io.Copy(w, content)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meta)
}

func (h *GcsHandler) handleDeleteObject(w http.ResponseWriter, r *http.Request, ctx context.Context, bucketName, objectName string) {
	h.service.DeleteObject(ctx, bucketName, objectName)
	w.WriteHeader(http.StatusNoContent)
}

func (h *GcsHandler) handleUploadObject(w http.ResponseWriter, r *http.Request, ctx context.Context, bucketName string) {
	objectName := r.URL.Query().Get("name")
	contentType := r.Header.Get("Content-Type")
	
	meta, err := h.service.InsertObject(ctx, bucketName, objectName, r.Body, contentType)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meta)
}
