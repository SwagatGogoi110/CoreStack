package cloudfunctions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/cloudfunctions/model"
)

type CloudFunctionsHandler struct {
	service *CloudFunctionsService
}

func NewCloudFunctionsHandler(service *CloudFunctionsService) *CloudFunctionsHandler {
	return &CloudFunctionsHandler{
		service: service,
	}
}

func (h *CloudFunctionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v2/")
	ctx := context.Background()

	if strings.Contains(path, "/functions") {
		if strings.HasSuffix(path, "/functions") {
			if r.Method == "GET" {
				h.handleListFunctions(w, r, ctx, path)
				return
			}
			if r.Method == "POST" {
				h.handleCreateFunction(w, r, ctx, path)
				return
			}
		}

		if r.Method == "GET" {
			h.handleGetFunction(w, r, ctx, path)
			return
		}
		if r.Method == "DELETE" {
			h.handleDeleteFunction(w, r, ctx, path)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "GCP Cloud Functions path not implemented: %s"}}`, path)
}

func (h *CloudFunctionsHandler) handleListFunctions(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	functions, _ := h.service.ListFunctions(ctx, parent)
	res := model.FunctionsList{
		Functions: functions,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *CloudFunctionsHandler) handleCreateFunction(w http.ResponseWriter, r *http.Request, ctx context.Context, path string) {
	var fn model.Function
	json.NewDecoder(r.Body).Decode(&fn)
	
	functionId := r.URL.Query().Get("functionId")
	name := path + "/" + functionId

	created, err := h.service.CreateFunction(ctx, name, &fn)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	
	// Cloud Functions creation usually returns an Operation
	// For simplicity, we'll just return the function for now or a dummy operation
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *CloudFunctionsHandler) handleGetFunction(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	fn, err := h.service.GetFunction(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fn)
}

func (h *CloudFunctionsHandler) handleDeleteFunction(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	h.service.DeleteFunction(ctx, name)
	w.WriteHeader(http.StatusNoContent)
}
