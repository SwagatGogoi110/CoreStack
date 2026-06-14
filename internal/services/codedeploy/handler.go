package codedeploy

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type CodeDeployJsonHandler struct {
	service *CodeDeployService
}

func NewCodeDeployJsonHandler(service *CodeDeployService) *CodeDeployJsonHandler {
	return &CodeDeployJsonHandler{
		service: service,
	}
}

func (h *CodeDeployJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	// Handle with and without prefix
	if action == "ListApplications" || action == "CodeDeploy_20141006.ListApplications" {
		return map[string]any{"applications": []any{}}, nil
	}
	return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
}

func (h *CodeDeployJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListApplications" || action == "CodeDeploy_20141006.ListApplications" {
		b := common.NewXmlBuilder()
		b.Start("ListApplicationsResponse").Start("ListApplicationsResult")
		b.Start("applications").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
