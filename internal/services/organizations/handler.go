package organizations

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type OrganizationsJsonHandler struct {
	service *OrganizationsService
}

func NewOrganizationsJsonHandler(service *OrganizationsService) *OrganizationsJsonHandler {
	return &OrganizationsJsonHandler{
		service: service,
	}
}

func (h *OrganizationsJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	switch action {
	case "DescribeOrganization":
		return map[string]any{"Organization": nil}, nil
	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *OrganizationsJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "DescribeOrganization" {
		b := common.NewXmlBuilder()
		// Namespace?
		b.Start("DescribeOrganizationResponse").Start("DescribeOrganizationResult")
		b.Start("Organization").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
