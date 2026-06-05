package bcmdataexports

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
	"github.com/hectorvent/cloudstack/internal/services/bcmdataexports/model"
)

type BcmDataExportsJsonHandler struct {
	service *BcmDataExportsService
}

func NewBcmDataExportsJsonHandler(service *BcmDataExportsService) *BcmDataExportsJsonHandler {
	return &BcmDataExportsJsonHandler{
		service: service,
	}
}

func (h *BcmDataExportsJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateExport":
		var req struct {
			Export *model.Export `json:"Export"`
		}
		json.Unmarshal(request, &req)
		export, err := h.service.CreateExport(ctx, req.Export)
		if err != nil {
			return nil, err
		}
		return map[string]string{"ExportArn": export.ExportArn}, nil

	case "GetExport":
		var req struct {
			ExportArn string `json:"ExportArn"`
		}
		json.Unmarshal(request, &req)
		export, err := h.service.GetExport(ctx, req.ExportArn)
		if err != nil {
			return nil, err
		}
		return map[string]any{"Export": export}, nil

	case "ListExports":
		exports, _ := h.service.ListExports(ctx)
		return map[string]any{"Exports": exports}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *BcmDataExportsJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("BcmDataExports does not support Query protocol")
}
