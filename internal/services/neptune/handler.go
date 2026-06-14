package neptune

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type NeptuneQueryHandler struct {
	service *NeptuneService
}

func NewNeptuneQueryHandler(service *NeptuneService) *NeptuneQueryHandler {
	return &NeptuneQueryHandler{
		service: service,
	}
}

func (h *NeptuneQueryHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	switch action {
	case "DescribeDBInstances":
		return map[string]any{"DBInstances": []any{}}, nil
	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *NeptuneQueryHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "DescribeDBInstances" {
		b := common.NewXmlBuilder()
		// Namespace?
		b.Start("DescribeDBInstancesResponse").Start("DescribeDBInstancesResult")
		b.Start("DBInstances").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
