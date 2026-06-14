package sagemaker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type SagemakerJsonHandler struct {
	service *SagemakerService
}

func NewSagemakerJsonHandler(service *SagemakerService) *SagemakerJsonHandler {
	return &SagemakerJsonHandler{
		service: service,
	}
}

func (h *SagemakerJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateNotebookInstance":
		var req struct {
			NotebookInstanceName string `json:"NotebookInstanceName"`
			InstanceType         string `json:"InstanceType"`
		}
		json.Unmarshal(request, &req)
		nb, err := h.service.CreateNotebookInstance(ctx, req.NotebookInstanceName, req.InstanceType)
		if err != nil {
			return nil, err
		}
		return map[string]string{"NotebookInstanceArn": nb.NotebookInstanceArn}, nil

	case "ListNotebookInstances":
		nbs, _ := h.service.ListNotebookInstances(ctx)
		return map[string]any{"NotebookInstances": nbs}, nil

	case "CreateModel":
		var req struct {
			ModelName string `json:"ModelName"`
		}
		json.Unmarshal(request, &req)
		m, err := h.service.CreateModel(ctx, req.ModelName)
		if err != nil {
			return nil, err
		}
		return map[string]string{"ModelArn": m.ModelArn}, nil

	case "ListModels":
		ms, _ := h.service.ListModels(ctx)
		return map[string]any{"Models": ms}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *SagemakerJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("SageMaker does not support Query protocol")
}
