package kinesis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type KinesisJsonHandler struct {
	service *KinesisService
}

func NewKinesisJsonHandler(service *KinesisService) *KinesisJsonHandler {
	return &KinesisJsonHandler{
		service: service,
	}
}

func (h *KinesisJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateStream":
		var req struct {
			StreamName string `json:"StreamName"`
			ShardCount int    `json:"ShardCount"`
		}
		json.Unmarshal(request, &req)
		if _, err := h.service.CreateStream(ctx, req.StreamName, req.ShardCount); err != nil {
			return nil, err
		}
		return map[string]any{}, nil

	case "ListStreams":
		names, _ := h.service.ListStreams(ctx)
		return map[string]any{"StreamNames": names, "HasMoreStreams": false}, nil

	case "PutRecord":
		var req struct {
			StreamName   string `json:"StreamName"`
			Data         []byte `json:"Data"`
			PartitionKey string `json:"PartitionKey"`
		}
		json.Unmarshal(request, &req)
		seq, err := h.service.PutRecord(ctx, req.StreamName, req.PartitionKey, req.Data)
		if err != nil {
			return nil, err
		}
		return map[string]any{"SequenceNumber": seq, "ShardId": "shardId-000000000000"}, nil

	default:
		return nil, fmt.Errorf("Unknown Kinesis action: %s", action)
	}
}

func (h *KinesisJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("Kinesis does not support Query protocol")
}
