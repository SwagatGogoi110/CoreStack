package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type BackupJsonHandler struct {
	service *BackupService
}

func NewBackupJsonHandler(service *BackupService) *BackupJsonHandler {
	return &BackupJsonHandler{
		service: service,
	}
}

func (h *BackupJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateBackupVault":
		var req struct {
			BackupVaultName string `json:"BackupVaultName"`
		}
		json.Unmarshal(request, &req)
		vault, err := h.service.CreateBackupVault(ctx, req.BackupVaultName)
		if err != nil {
			return nil, err
		}
		return vault, nil

	case "DescribeBackupVault":
		var req struct {
			BackupVaultName string `json:"BackupVaultName"`
		}
		json.Unmarshal(request, &req)
		vault, err := h.service.DescribeBackupVault(ctx, req.BackupVaultName)
		if err != nil {
			return nil, err
		}
		return vault, nil

	case "ListBackupVaults":
		vaults, _ := h.service.ListBackupVaults(ctx)
		return map[string]any{"BackupVaultList": vaults}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *BackupJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("Backup does not support Query protocol")
}
