package eks

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type EksJsonHandler struct {
	service *EksService
}

func NewEksJsonHandler(service *EksService) *EksJsonHandler {
	return &EksJsonHandler{
		service: service,
	}
}

func (h *EksJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	// Handle with and without prefix
	if action == "ListClusters" || action == "AWSClusterManagementService.ListClusters" {
		return map[string]any{"clusters": []string{}}, nil
	}
	return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
}

func (h *EksJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListClusters" || action == "AWSClusterManagementService.ListClusters" {
		b := common.NewXmlBuilder()
		b.Start("ListClustersResponse").Start("ListClustersResult")
		b.Start("clusters").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
