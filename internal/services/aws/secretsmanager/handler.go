package secretsmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type SecretsManagerJsonHandler struct {
	service *SecretsManagerService
}

func NewSecretsManagerJsonHandler(service *SecretsManagerService) *SecretsManagerJsonHandler {
	return &SecretsManagerJsonHandler{
		service: service,
	}
}

func (h *SecretsManagerJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateSecret":
		var req struct {
			Name         string `json:"Name"`
			SecretString string `json:"SecretString"`
		}
		json.Unmarshal(request, &req)
		secret, err := h.service.CreateSecret(ctx, req.Name, req.SecretString)
		if err != nil {
			return nil, err
		}
		return map[string]string{"ARN": secret.Arn, "Name": secret.Name}, nil

	case "GetSecretValue":
		var req struct {
			SecretID string `json:"SecretId"`
		}
		json.Unmarshal(request, &req)
		val, err := h.service.GetSecretValue(ctx, req.SecretID)
		if err != nil {
			return nil, err
		}
		return val, nil

	case "ListSecrets":
		secrets, _ := h.service.ListSecrets(ctx)
		return map[string]any{"SecretList": secrets}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *SecretsManagerJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("SecretsManager does not support Query protocol")
}
