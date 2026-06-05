package transcribe

import (
	"context"
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
	ctx := context.Background()

	switch action {
	case "StartTranscriptionJob":
		var req struct {
			TranscriptionJobName string `json:"TranscriptionJobName"`
			LanguageCode         string `json:"LanguageCode"`
			MediaFormat          string `json:"MediaFormat"`
		}
		json.Unmarshal(request, &req)
		job, err := h.service.StartTranscriptionJob(ctx, req.TranscriptionJobName, req.LanguageCode, req.MediaFormat)
		if err != nil {
			return nil, err
		}
		return map[string]any{"TranscriptionJob": job}, nil

	case "GetTranscriptionJob":
		var req struct {
			TranscriptionJobName string `json:"TranscriptionJobName"`
		}
		json.Unmarshal(request, &req)
		job, err := h.service.GetTranscriptionJob(ctx, req.TranscriptionJobName)
		if err != nil {
			return nil, err
		}
		return map[string]any{"TranscriptionJob": job}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *TranscribeJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("Transcribe does not support Query protocol")
}
