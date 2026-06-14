package ecr

import (
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
	switch action {
	case "DescribeRepositories":
		return map[string]any{"repositories": []any{}}, nil
	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *EcrJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "DescribeRepositories" {
		b := common.NewXmlBuilder()
		// Namespace?
		b.Start("DescribeRepositoriesResponse").Start("DescribeRepositoriesResult")
		b.Start("repositories").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
