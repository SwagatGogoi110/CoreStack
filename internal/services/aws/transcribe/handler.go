package transcribe

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type TranscribeJsonHandler struct {
	service *TranscribeService
}

func NewTranscribeJsonHandler(service *TranscribeService) *TranscribeJsonHandler {
	return &TranscribeJsonHandler{
		service: service,
	}
}

func (h *TranscribeJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	switch action {
	case "ListTranscriptionJobs":
		return map[string]any{"TranscriptionJobSummaries": []any{}}, nil
	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *TranscribeJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListTranscriptionJobs" {
		b := common.NewXmlBuilder()
		// Namespace?
		b.Start("ListTranscriptionJobsResponse").Start("ListTranscriptionJobsResult")
		b.Start("TranscriptionJobSummaries").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
