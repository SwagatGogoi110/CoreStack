package kms

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type KmsJsonHandler struct {
	service *KmsService
}

func NewKmsJsonHandler(service *KmsService) *KmsJsonHandler {
	return &KmsJsonHandler{
		service: service,
	}
}

func (h *KmsJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	switch action {
	case "ListKeys":
		return map[string]any{"Keys": []any{}}, nil
	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *KmsJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListKeys" {
		b := common.NewXmlBuilder()
		// Namespace?
		b.Start("ListKeysResponse").Start("ListKeysResult")
		b.Start("Keys").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
