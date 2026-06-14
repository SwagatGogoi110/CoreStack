package workflows

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/workflows/model"
)

type WorkflowsHandler struct {
	service *WorkflowsService
}

func NewWorkflowsHandler(service *WorkflowsService) *WorkflowsHandler {
	return &WorkflowsHandler{
		service: service,
	}
}

func (h *WorkflowsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/")
	ctx := context.Background()

	if strings.Contains(path, "/workflows") {
		parts := strings.Split(path, "/")
		
		if strings.HasSuffix(path, "/workflows") {
			if r.Method == "GET" {
				h.handleListWorkflows(w, r, ctx, path)
				return
			}
			if r.Method == "POST" {
				h.handleCreateWorkflow(w, r, ctx, path)
				return
			}
		}

		if strings.Contains(path, "/executions") {
			idx := -1
			for i, p := range parts {
				if p == "executions" {
					idx = i
					break
				}
			}
			
			workflowName := strings.Join(parts[:idx], "/")
			if strings.HasSuffix(path, "/executions") {
				if r.Method == "GET" {
					h.handleListExecutions(w, r, ctx, workflowName)
					return
				}
				if r.Method == "POST" {
					h.handleCreateExecution(w, r, ctx, workflowName)
					return
				}
			}
		}

		if r.Method == "GET" {
			h.handleGetWorkflow(w, r, ctx, path)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "Workflows path not implemented: %s"}}`, path)
}

func (h *WorkflowsHandler) handleListWorkflows(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	items, _ := h.service.ListWorkflows(ctx, parent)
	res := model.WorkflowsList{
		Workflows: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *WorkflowsHandler) handleCreateWorkflow(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	var wf model.Workflow
	json.NewDecoder(r.Body).Decode(&wf)
	
	workflowId := r.URL.Query().Get("workflowId")
	name := parent + "/" + workflowId

	created, _ := h.service.CreateWorkflow(ctx, name, &wf)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *WorkflowsHandler) handleGetWorkflow(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	items, _ := h.service.ListWorkflows(ctx, "")
	for _, wf := range items {
		if wf.Name == name {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(wf)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func (h *WorkflowsHandler) handleListExecutions(w http.ResponseWriter, r *http.Request, ctx context.Context, workflowName string) {
	items, _ := h.service.ListExecutions(ctx, workflowName)
	res := model.ExecutionsList{
		Executions: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *WorkflowsHandler) handleCreateExecution(w http.ResponseWriter, r *http.Request, ctx context.Context, workflowName string) {
	var exec model.Execution
	json.NewDecoder(r.Body).Decode(&exec)
	created, _ := h.service.CreateExecution(ctx, workflowName, &exec)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}
