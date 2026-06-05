package kms

import (
	"context"
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
	ctx := context.Background()

	switch action {
	case "CreateKey":
		var req struct {
			Description string `json:"Description"`
		}
		json.Unmarshal(request, &req)
		key, err := h.service.CreateKey(ctx, req.Description)
		if err != nil {
			return nil, err
		}
		return map[string]any{"KeyMetadata": key}, nil

	case "DescribeKey":
		var req struct {
			KeyID string `json:"KeyId"`
		}
		json.Unmarshal(request, &req)
		key, err := h.service.DescribeKey(ctx, req.KeyID)
		if err != nil {
			return nil, err
		}
		return map[string]any{"KeyMetadata": key}, nil

	case "Encrypt":
		var req struct {
			KeyID     string `json:"KeyId"`
			Plaintext []byte `json:"Plaintext"`
		}
		json.Unmarshal(request, &req)
		cipher, err := h.service.Encrypt(ctx, req.KeyID, req.Plaintext)
		if err != nil {
			return nil, err
		}
		return map[string]any{"CiphertextBlob": cipher, "KeyId": req.KeyID}, nil

	case "Decrypt":
		var req struct {
			CiphertextBlob []byte `json:"CiphertextBlob"`
		}
		json.Unmarshal(request, &req)
		plain, err := h.service.Decrypt(ctx, "", req.CiphertextBlob)
		if err != nil {
			return nil, err
		}
		return map[string]any{"Plaintext": plain, "KeyId": "stub-key-id"}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *KmsJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("KMS does not support Query protocol")
}
