package codecommit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type CodecommitJsonHandler struct {
	service *CodecommitService
}

func NewCodecommitJsonHandler(service *CodecommitService) *CodecommitJsonHandler {
	return &CodecommitJsonHandler{
		service: service,
	}
}

func (h *CodecommitJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateRepository":
		var req struct {
			RepositoryName        string `json:"repositoryName"`
			RepositoryDescription string `json:"repositoryDescription"`
		}
		json.Unmarshal(request, &req)
		repo, err := h.service.CreateRepository(ctx, req.RepositoryName, req.RepositoryDescription)
		if err != nil {
			return nil, err
		}
		return map[string]any{"repositoryMetadata": repo}, nil

	case "GetRepository":
		var req struct {
			RepositoryName string `json:"repositoryName"`
		}
		json.Unmarshal(request, &req)
		repo, err := h.service.GetRepository(ctx, req.RepositoryName)
		if err != nil {
			return nil, err
		}
		return map[string]any{"repositoryMetadata": repo}, nil

	case "ListRepositories":
		repos, _ := h.service.ListRepositories(ctx)
		return map[string]any{"repositories": repos}, nil

	case "DeleteRepository":
		var req struct {
			RepositoryName string `json:"repositoryName"`
		}
		json.Unmarshal(request, &req)
		h.service.DeleteRepository(ctx, req.RepositoryName)
		return map[string]any{}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *CodecommitJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("CodeCommit does not support Query protocol")
}
