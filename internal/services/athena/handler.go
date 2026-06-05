package athena

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type AthenaJsonHandler struct {
	service *AthenaService
}

func NewAthenaJsonHandler(service *AthenaService) *AthenaJsonHandler {
	return &AthenaJsonHandler{
		service: service,
	}
}

func (h *AthenaJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "StartQueryExecution":
		var req struct {
			Query     string `json:"QueryString"`
			WorkGroup string `json:"WorkGroup"`
		}
		json.Unmarshal(request, &req)
		id, err := h.service.StartQueryExecution(ctx, req.Query, req.WorkGroup)
		if err != nil {
			return nil, err
		}
		return map[string]string{"QueryExecutionId": id}, nil

	case "GetQueryExecution":
		var req struct {
			QueryExecutionID string `json:"QueryExecutionId"`
		}
		json.Unmarshal(request, &req)
		execution, err := h.service.GetQueryExecution(ctx, req.QueryExecutionID)
		if err != nil {
			return nil, err
		}
		return map[string]any{"QueryExecution": execution}, nil

	case "ListQueryExecutions":
		return map[string]any{"QueryExecutionIds": []string{}}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *AthenaJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("Athena does not support Query protocol")
}
