package operations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/operations/model"
)

type OperationsHandler struct {
	service *OperationsService
}

func NewOperationsHandler(service *OperationsService) *OperationsHandler {
	return &OperationsHandler{
		service: service,
	}
}

func (h *OperationsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/")
	ctx := context.Background()

	if strings.Contains(path, "/operations") {
		if strings.HasSuffix(path, "/operations") {
			if r.Method == "GET" {
				h.handleListOperations(w, r, ctx, path)
				return
			}
		}

		if r.Method == "GET" {
			h.handleGetOperation(w, r, ctx, path)
			return
		}
		if r.Method == "DELETE" {
			h.handleDeleteOperation(w, r, ctx, path)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "GCP Operations path not implemented: %s"}}`, path)
}

func (h *OperationsHandler) handleListOperations(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	ops, _ := h.service.ListOperations(ctx, name)
	res := model.OperationsList{
		Operations: ops,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *OperationsHandler) handleGetOperation(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	op, err := h.service.GetOperation(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(op)
}

func (h *OperationsHandler) handleDeleteOperation(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	h.service.DeleteOperation(ctx, name)
	w.WriteHeader(http.StatusNoContent)
}
