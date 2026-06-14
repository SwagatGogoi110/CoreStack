package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/tasks/model"
)

type TasksHandler struct {
	service *TasksService
}

func NewTasksHandler(service *TasksService) *TasksHandler {
	return &TasksHandler{
		service: service,
	}
}

func (h *TasksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v2/")
	ctx := context.Background()

	if strings.Contains(path, "/queues") {
		if strings.HasSuffix(path, "/queues") {
			if r.Method == "GET" {
				h.handleListQueues(w, r, ctx, path)
				return
			}
			if r.Method == "POST" {
				h.handleCreateQueue(w, r, ctx, path)
				return
			}
		}

		if strings.Contains(path, "/tasks") {
			if strings.HasSuffix(path, "/tasks") {
				if r.Method == "GET" {
					h.handleListTasks(w, r, ctx, path)
					return
				}
				if r.Method == "POST" {
					h.handleCreateTask(w, r, ctx, path)
					return
				}
			}
		}

		// Delete
		if r.Method == "DELETE" {
			h.handleDeleteQueue(w, r, ctx, path)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "GCP Cloud Tasks path not implemented: %s"}}`, path)
}

func (h *TasksHandler) handleListQueues(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	queues, _ := h.service.ListQueues(ctx, parent)
	res := model.QueuesList{
		Queues: queues,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *TasksHandler) handleCreateQueue(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	var q model.Queue
	json.NewDecoder(r.Body).Decode(&q)
	
	// Create full name if not provided
	if q.Name == "" {
		// Mock name generation
		q.Name = parent + "/my-queue"
	}

	created, _ := h.service.CreateQueue(ctx, q.Name, &q)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *TasksHandler) handleListTasks(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	tasks, _ := h.service.ListTasks(ctx, parent)
	res := model.TasksList{
		Tasks: tasks,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *TasksHandler) handleCreateTask(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	var task model.Task
	json.NewDecoder(r.Body).Decode(&task)
	
	if task.Name == "" {
		task.Name = parent + "/tasks/task1"
	}

	created, _ := h.service.CreateTask(ctx, parent, &task)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *TasksHandler) handleDeleteQueue(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	h.service.DeleteQueue(ctx, name)
	w.WriteHeader(http.StatusNoContent)
}
