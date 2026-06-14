package firehose

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
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
	switch action {
	case "ListDeliveryStreams":
		return map[string]any{"DeliveryStreamNames": []string{}}, nil
	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *FirehoseJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListDeliveryStreams" {
		b := common.NewXmlBuilder()
		// Namespace?
		b.Start("ListDeliveryStreamsResponse").Start("ListDeliveryStreamsResult")
		b.Start("DeliveryStreamNames").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
