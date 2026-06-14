package cur

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
	"github.com/hectorvent/cloudstack/internal/services/aws/cur/model"
)

type CurJsonHandler struct {
	service *CurService
}

func NewCurJsonHandler(service *CurService) *CurJsonHandler {
	return &CurJsonHandler{
		service: service,
	}
}

func (h *CurJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "PutReportDefinition":
		var req struct {
			ReportDefinition *model.ReportDefinition `json:"ReportDefinition"`
		}
		json.Unmarshal(request, &req)
		if err := h.service.PutReportDefinition(ctx, req.ReportDefinition); err != nil {
			return nil, err
		}
		return map[string]any{}, nil

	case "DescribeReportDefinitions":
		reports, _ := h.service.DescribeReportDefinitions(ctx)
		return map[string]any{"ReportDefinitions": reports}, nil

	case "DeleteReportDefinition":
		var req struct {
			ReportName string `json:"ReportName"`
		}
		json.Unmarshal(request, &req)
		// ...
		return map[string]any{}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *CurJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("CUR does not support Query protocol")
}
