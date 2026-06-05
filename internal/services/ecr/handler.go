package ecr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type EcrJsonHandler struct {
	service *EcrService
}

func NewEcrJsonHandler(service *EcrService) *EcrJsonHandler {
	return &EcrJsonHandler{
		service: service,
	}
}

func (h *EcrJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "GetAuthorizationToken":
		return map[string]any{
			"authorizationData": []any{
				map[string]any{
					"authorizationToken": "QVdTOmZsb2Np", // AWS:CloudStack
					"proxyEndpoint":      "http://localhost:8080",
				},
			},
		}, nil

	case "CreateRepository":
		var req struct {
			RepositoryName string `json:"repositoryName"`
		}
		json.Unmarshal(request, &req)
		repo, err := h.service.CreateRepository(ctx, req.RepositoryName)
		if err != nil {
			return nil, err
		}
		return map[string]any{"repository": repo}, nil

	case "DescribeRepositories":
		var req struct {
			RepositoryNames []string `json:"repositoryNames"`
		}
		json.Unmarshal(request, &req)
		repos, _ := h.service.DescribeRepositories(ctx, req.RepositoryNames)
		return map[string]any{"repositories": repos}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *EcrJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("ECR does not support Query protocol")
}
