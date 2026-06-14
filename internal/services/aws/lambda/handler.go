package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hectorvent/cloudstack/internal/core/common"
	"github.com/hectorvent/cloudstack/internal/services/aws/lambda/model"
)

type LambdaHandler struct {
	service *LambdaService
}

func NewLambdaHandler(service *LambdaService) *LambdaHandler {
	return &LambdaHandler{
		service: service,
	}
}

func (h *LambdaHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("Lambda does not support standard JSON 1.0/1.1 protocol dispatcher")
}

func (h *LambdaHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("Lambda does not support Query protocol")
}

func (h *LambdaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/2015-03-31")
	
	if path == "/functions" && r.Method == "POST" {
		h.handleCreateFunction(w, r)
		return
	}
	if path == "/functions" && r.Method == "GET" {
		h.handleListFunctions(w, r)
		return
	}
	
	if strings.HasPrefix(path, "/functions/") {
		parts := strings.Split(strings.TrimPrefix(path, "/functions/"), "/")
		functionName := parts[0]
		
		if len(parts) == 1 && r.Method == "GET" {
			h.handleGetFunction(w, r, functionName)
			return
		}
		if len(parts) == 2 && parts[1] == "invocations" && r.Method == "POST" {
			h.handleInvoke(w, r, functionName)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func (h *LambdaHandler) handleCreateFunction(w http.ResponseWriter, r *http.Request) {
	var body struct {
		FunctionName string            `json:"FunctionName"`
		Runtime      string            `json:"Runtime"`
		Role         string            `json:"Role"`
		Handler      string            `json:"Handler"`
		Code         struct {
			ZipFile []byte `json:"ZipFile"`
		} `json:"Code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fn := &model.LambdaFunction{
		FunctionName: body.FunctionName,
		Runtime:      body.Runtime,
		Role:         body.Role,
		Handler:      body.Handler,
	}

	ctx := context.Background()
	created, err := h.service.CreateFunction(ctx, fn, body.Code.ZipFile)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *LambdaHandler) handleGetFunction(w http.ResponseWriter, r *http.Request, functionName string) {
	ctx := context.Background()
	fn, err := h.service.GetFunction(ctx, functionName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res := map[string]any{
		"Configuration": fn,
		"Code": map[string]string{
			"Location": "http://localhost:8080/fake-s3-link",
		},
	}
	json.NewEncoder(w).Encode(res)
}

func (h *LambdaHandler) handleListFunctions(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	fns, _ := h.service.ListFunctions(ctx)
	res := map[string]any{
		"Functions": fns,
	}
	json.NewEncoder(w).Encode(res)
}

func (h *LambdaHandler) handleInvoke(w http.ResponseWriter, r *http.Request, functionName string) {
	payload, _ := io.ReadAll(r.Body)
	ctx := context.Background()
	result, err := h.service.Invoke(ctx, functionName, payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result.Payload)
}
