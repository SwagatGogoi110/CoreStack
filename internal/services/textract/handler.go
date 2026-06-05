package textract

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type TextractJsonHandler struct {
	service *TextractService
}

func NewTextractJsonHandler(service *TextractService) *TextractJsonHandler {
	return &TextractJsonHandler{
		service: service,
	}
}

func (h *TextractJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	switch action {
	case "DetectDocumentText":
		return h.service.DetectDocumentText(), nil
	case "AnalyzeDocument":
		return h.service.DetectDocumentText(), nil
	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *TextractJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("Textract does not support Query protocol")
}
