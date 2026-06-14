package codepipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
	"github.com/hectorvent/cloudstack/internal/services/codepipeline/model"
)

type CodepipelineJsonHandler struct {
	service *CodepipelineService
}

func NewCodepipelineJsonHandler(service *CodepipelineService) *CodepipelineJsonHandler {
	return &CodepipelineJsonHandler{
		service: service,
	}
}

func (h *CodepipelineJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreatePipeline":
		var req struct {
			Pipeline *model.Pipeline `json:"pipeline"`
		}
		json.Unmarshal(request, &req)
		p, err := h.service.CreatePipeline(ctx, req.Pipeline)
		if err != nil {
			return nil, err
		}
		return map[string]any{"pipeline": p}, nil

	case "GetPipeline":
		var req struct {
			Name string `json:"name"`
		}
		json.Unmarshal(request, &req)
		p, err := h.service.GetPipeline(ctx, req.Name)
		if err != nil {
			return nil, err
		}
		return map[string]any{"pipeline": p}, nil

	case "ListPipelines":
		summaries, _ := h.service.ListPipelines(ctx)
		return map[string]any{"pipelines": summaries}, nil

	case "DeletePipeline":
		var req struct {
			Name string `json:"name"`
		}
		json.Unmarshal(request, &req)
		h.service.DeletePipeline(ctx, req.Name)
		return map[string]any{}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *CodepipelineJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("CodePipeline does not support Query protocol")
}
