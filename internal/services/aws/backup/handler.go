package backup

import (
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
	// Handle with and without prefix
	if action == "ListBackupPlans" || action == "BackupStorage.ListBackupPlans" {
		return map[string]any{"BackupPlansList": []any{}}, nil
	}
	return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
}

func (h *BackupJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListBackupPlans" || action == "BackupStorage.ListBackupPlans" {
		b := common.NewXmlBuilder()
		b.Start("ListBackupPlansResponse").Start("ListBackupPlansResult")
		b.Start("BackupPlansList").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
