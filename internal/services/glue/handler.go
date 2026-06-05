package glue

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
	"github.com/hectorvent/cloudstack/internal/services/glue/model"
)

type GlueJsonHandler struct {
	service *GlueService
}

func NewGlueJsonHandler(service *GlueService) *GlueJsonHandler {
	return &GlueJsonHandler{
		service: service,
	}
}

func (h *GlueJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateDatabase":
		var req struct {
			DatabaseInput *model.Database `json:"DatabaseInput"`
		}
		json.Unmarshal(request, &req)
		db, err := h.service.CreateDatabase(ctx, req.DatabaseInput.Name, req.DatabaseInput.Description)
		if err != nil {
			return nil, err
		}
		return db, nil

	case "GetDatabase":
		var req struct {
			Name string `json:"Name"`
		}
		json.Unmarshal(request, &req)
		db, err := h.service.GetDatabase(ctx, req.Name)
		if err != nil {
			return nil, err
		}
		return map[string]any{"Database": db}, nil

	case "CreateTable":
		var req struct {
			DatabaseName string       `json:"DatabaseName"`
			TableInput   *model.Table `json:"TableInput"`
		}
		json.Unmarshal(request, &req)
		if err := h.service.CreateTable(ctx, req.DatabaseName, req.TableInput); err != nil {
			return nil, err
		}
		return map[string]any{}, nil

	case "GetTables":
		var req struct {
			DatabaseName string `json:"DatabaseName"`
		}
		json.Unmarshal(request, &req)
		tables, _ := h.service.GetTables(ctx, req.DatabaseName)
		return map[string]any{"TableList": tables}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *GlueJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("Glue does not support Query protocol")
}
