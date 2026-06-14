package dynamodb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
	"github.com/hectorvent/cloudstack/internal/services/dynamodb/model"
)

type DynamoDbHandler struct {
	service *DynamoDbService
}

func NewDynamoDbHandler(service *DynamoDbService) *DynamoDbHandler {
	return &DynamoDbHandler{
		service: service,
	}
}

func (h *DynamoDbHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateTable":
		var req model.TableDefinition
		if err := json.Unmarshal(request, &req); err != nil {
			return nil, err
		}
		table, err := h.service.CreateTable(ctx, &req)
		if err != nil {
			return nil, err
		}
		return map[string]any{"TableDescription": table}, nil

	case "DescribeTable":
		var req struct {
			TableName string `json:"TableName"`
		}
		if err := json.Unmarshal(request, &req); err != nil {
			return nil, err
		}
		table, err := h.service.DescribeTable(ctx, req.TableName)
		if err != nil {
			return nil, err
		}
		return map[string]any{"Table": table}, nil

	case "ListTables":
		tables, _ := h.service.ListTables(ctx)
		return map[string]any{"TableNames": tables}, nil

	case "PutItem":
		var req struct {
			TableName string          `json:"TableName"`
			Item      json.RawMessage `json:"Item"`
		}
		if err := json.Unmarshal(request, &req); err != nil {
			return nil, err
		}
		if err := h.service.PutItem(ctx, req.TableName, req.Item); err != nil {
			return nil, err
		}
		return map[string]any{}, nil

	case "GetItem":
		var req struct {
			TableName string          `json:"TableName"`
			Key       json.RawMessage `json:"Key"`
		}
		if err := json.Unmarshal(request, &req); err != nil {
			return nil, err
		}
		item, err := h.service.GetItem(ctx, req.TableName, req.Key)
		if err != nil {
			return nil, err
		}
		if item == nil {
			return map[string]any{}, nil
		}
		var m map[string]any
		json.Unmarshal(item, &m)
		return map[string]any{"Item": m}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *DynamoDbHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("DynamoDB does not support Query protocol")
}
