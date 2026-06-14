package bcmdataexports

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type BcmDataExportsJsonHandler struct {
	service *BcmDataExportsService
}

func NewBcmDataExportsJsonHandler(service *BcmDataExportsService) *BcmDataExportsJsonHandler {
	return &BcmDataExportsJsonHandler{
		service: service,
	}
}

func (h *BcmDataExportsJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	// Handle with and without prefix
	if action == "ListExports" || action == "BillingAndCostManagementDataExports.ListExports" {
		return map[string]any{"Exports": []any{}}, nil
	}
	return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
}

func (h *BcmDataExportsJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListExports" || action == "BillingAndCostManagementDataExports.ListExports" {
		b := common.NewXmlBuilder()
		b.Start("ListExportsResponse").Start("ListExportsResult")
		b.Start("Exports").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
