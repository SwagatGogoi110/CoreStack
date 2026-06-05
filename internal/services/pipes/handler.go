package pipes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type PipesJsonHandler struct {
	service *PipesService
}

func NewPipesJsonHandler(service *PipesService) *PipesJsonHandler {
	return &PipesJsonHandler{
		service: service,
	}
}

func (h *PipesJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreatePipe":
		var req struct {
			Name   string `json:"Name"`
			Source string `json:"Source"`
			Target string `json:"Target"`
		}
		json.Unmarshal(request, &req)
		pipe, err := h.service.CreatePipe(ctx, req.Name, req.Source, req.Target)
		if err != nil {
			return nil, err
		}
		return map[string]string{"Arn": pipe.Arn, "Name": pipe.Name}, nil

	case "DescribePipe":
		var req struct {
			Name string `json:"Name"`
		}
		json.Unmarshal(request, &req)
		pipe, err := h.service.DescribePipe(ctx, req.Name)
		if err != nil {
			return nil, err
		}
		return pipe, nil

	case "ListPipes":
		pipes, _ := h.service.ListPipes(ctx)
		return map[string]any{"Pipes": pipes}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *PipesJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("Pipes does not support Query protocol")
}
