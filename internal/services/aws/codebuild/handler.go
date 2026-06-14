package codebuild

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type CodeBuildJsonHandler struct {
	service *CodeBuildService
}

func NewCodeBuildJsonHandler(service *CodeBuildService) *CodeBuildJsonHandler {
	return &CodeBuildJsonHandler{
		service: service,
	}
}

func (h *CodeBuildJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateProject":
		var req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		json.Unmarshal(request, &req)
		project, err := h.service.CreateProject(ctx, req.Name, req.Description)
		if err != nil {
			return nil, err
		}
		return map[string]any{"project": project}, nil

	case "BatchGetProjects":
		var req struct {
			Names []string `json:"names"`
		}
		json.Unmarshal(request, &req)
		projects, _ := h.service.BatchGetProjects(ctx, req.Names)
		return map[string]any{"projects": projects}, nil

	case "ListProjects":
		names, _ := h.service.ListProjects(ctx)
		return map[string]any{"projects": names}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *CodeBuildJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("CodeBuild does not support Query protocol")
}
