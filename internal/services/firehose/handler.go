package firehose

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
	"github.com/hectorvent/cloudstack/internal/services/firehose/model"
)

type FirehoseJsonHandler struct {
	service *FirehoseService
}

func NewFirehoseJsonHandler(service *FirehoseService) *FirehoseJsonHandler {
	return &FirehoseJsonHandler{
		service: service,
	}
}

func (h *FirehoseJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateDeliveryStream":
		var req struct {
			DeliveryStreamName string `json:"DeliveryStreamName"`
			S3Destination      *model.S3DestinationDescription `json:"ExtendedS3DestinationConfiguration"`
		}
		json.Unmarshal(request, &req)
		stream, err := h.service.CreateDeliveryStream(ctx, req.DeliveryStreamName, req.S3Destination)
		if err != nil {
			return nil, err
		}
		return map[string]string{"DeliveryStreamARN": stream.DeliveryStreamArn}, nil

	case "PutRecord":
		var req struct {
			DeliveryStreamName string       `json:"DeliveryStreamName"`
			Record             model.Record `json:"Record"`
		}
		json.Unmarshal(request, &req)
		if err := h.service.PutRecord(ctx, req.DeliveryStreamName, req.Record.Data); err != nil {
			return nil, err
		}
		return map[string]string{"RecordId": "CloudStack-fh-rec"}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *FirehoseJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("Firehose does not support Query protocol")
}
