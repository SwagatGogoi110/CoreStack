package s3

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type S3Handler struct {
	service *S3Service
}

func NewS3Handler(service *S3Service) *S3Handler {
	return &S3Handler{
		service: service,
	}
}

func (h *S3Handler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("S3 does not support JSON protocol")
}

func (h *S3Handler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("S3 does not support Query protocol")
}

func (h *S3Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// S3 routing logic
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" && r.Method == "GET" {
		h.handleListBuckets(w, r)
		return
	}

	parts := strings.SplitN(path, "/", 2)
	bucket := parts[0]
	key := ""
	if len(parts) > 1 {
		key = parts[1]
	}

	if key == "" {
		// Bucket operations
		switch r.Method {
		case "PUT":
			h.handleCreateBucket(w, r, bucket)
		case "GET":
			h.handleListObjects(w, r, bucket)
		case "DELETE":
			h.handleDeleteBucket(w, r, bucket)
		case "HEAD":
			h.handleHeadBucket(w, r, bucket)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else {
		// Object operations
		switch r.Method {
		case "PUT":
			h.handlePutObject(w, r, bucket, key)
		case "GET":
			h.handleGetObject(w, r, bucket, key)
		case "DELETE":
			h.handleDeleteObject(w, r, bucket, key)
		case "HEAD":
			h.handleHeadObject(w, r, bucket, key)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (h *S3Handler) handleListBuckets(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	buckets, err := h.service.ListBuckets(ctx)
	if err != nil {
		h.writeError(w, "InternalError", err.Error(), http.StatusInternalServerError)
		return
	}

	xml := common.NewXmlBuilder().
		Start("ListAllMyBucketsResult").
		Start("Owner").
		Elem("ID", "CloudStack-owner").
		Elem("DisplayName", "CloudStack").
		End().
		Start("Buckets")

	for _, b := range buckets {
		xml.Start("Bucket").
			Elem("Name", b.Name).
			Elem("CreationDate", b.CreationDate.Format("2006-01-02T15:04:05.000Z")).
			End()
	}

	xml.End().End()

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, xml.Build())
}

func (h *S3Handler) handleCreateBucket(w http.ResponseWriter, r *http.Request, bucket string) {
	ctx := context.Background()
	_, err := h.service.CreateBucket(ctx, bucket, "us-east-1")
	if err != nil {
		h.writeError(w, "BucketAlreadyExists", err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *S3Handler) handlePutObject(w http.ResponseWriter, r *http.Request, bucket, key string) {
	ctx := context.Background()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, "InternalError", err.Error(), http.StatusInternalServerError)
		return
	}

	obj, err := h.service.PutObject(ctx, bucket, key, data, r.Header.Get("Content-Type"), nil)
	if err != nil {
		h.writeError(w, "InternalError", err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("ETag", obj.ETag)
	w.WriteHeader(http.StatusOK)
}

func (h *S3Handler) handleGetObject(w http.ResponseWriter, r *http.Request, bucket, key string) {
	ctx := context.Background()
	obj, data, err := h.service.GetObject(ctx, bucket, key)
	if err != nil {
		h.writeError(w, "NoSuchKey", err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", obj.ContentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", obj.Size))
	w.Header().Set("ETag", obj.ETag)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Stubs for now
func (h *S3Handler) handleListObjects(w http.ResponseWriter, r *http.Request, bucket string) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (h *S3Handler) handleDeleteBucket(w http.ResponseWriter, r *http.Request, bucket string) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (h *S3Handler) handleHeadBucket(w http.ResponseWriter, r *http.Request, bucket string) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (h *S3Handler) handleDeleteObject(w http.ResponseWriter, r *http.Request, bucket, key string) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (h *S3Handler) handleHeadObject(w http.ResponseWriter, r *http.Request, bucket, key string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *S3Handler) writeError(w http.ResponseWriter, code, message string, status int) {
	xml := common.NewXmlBuilder().
		Start("Error").
		Elem("Code", code).
		Elem("Message", message).
		Elem("RequestId", "CloudStack-s3-request").
		End().
		Build()
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	fmt.Fprint(w, xml)
}
