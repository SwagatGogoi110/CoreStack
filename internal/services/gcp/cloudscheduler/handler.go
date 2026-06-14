package cloudscheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/cloudscheduler/model"
)

type CloudSchedulerHandler struct {
	service *CloudSchedulerService
}

func NewCloudSchedulerHandler(service *CloudSchedulerService) *CloudSchedulerHandler {
	return &CloudSchedulerHandler{
		service: service,
	}
}

func (h *CloudSchedulerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/")
	ctx := context.Background()

	if strings.Contains(path, "/jobs") {
		
		if strings.HasSuffix(path, "/jobs") {
			if r.Method == "GET" {
				h.handleListJobs(w, r, ctx, path)
				return
			}
			if r.Method == "POST" {
				h.handleCreateJob(w, r, ctx, path)
				return
			}
		}

		if strings.Contains(path, ":run") {
			jobName := strings.TrimSuffix(path, ":run")
			h.handleRunJob(w, r, ctx, jobName)
			return
		}

		if r.Method == "GET" {
			h.handleGetJob(w, r, ctx, path)
			return
		}
		if r.Method == "DELETE" {
			h.handleDeleteJob(w, r, ctx, path)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "Cloud Scheduler path not implemented: %s"}}`, path)
}

func (h *CloudSchedulerHandler) handleListJobs(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	items, _ := h.service.ListJobs(ctx, parent)
	res := model.JobsList{
		Jobs: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *CloudSchedulerHandler) handleCreateJob(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	var job model.Job
	json.NewDecoder(r.Body).Decode(&job)
	
	// Create full name if not provided
	if job.Name == "" {
		// Mock name generation from request would be better
		job.Name = parent + "/my-job"
	}

	created, _ := h.service.CreateJob(ctx, job.Name, &job)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *CloudSchedulerHandler) handleGetJob(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	items, _ := h.service.ListJobs(ctx, "")
	for _, job := range items {
		if job.Name == name {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(job)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func (h *CloudSchedulerHandler) handleRunJob(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	job, err := h.service.RunJob(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (h *CloudSchedulerHandler) handleDeleteJob(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	h.service.DeleteJob(ctx, name)
	w.WriteHeader(http.StatusNoContent)
}
