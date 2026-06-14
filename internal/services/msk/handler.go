package msk

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type MskJsonHandler struct {
	service *MskService
}

func NewMskJsonHandler(service *MskService) *MskJsonHandler {
	return &MskJsonHandler{
		service: service,
	}
}

func (h *MskJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	// Handle with and without prefix
	if action == "ListClusters" || action == "AmazonMSK.ListClusters" {
		return map[string]any{"clusters": []any{}}, nil
	}
	return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
}

func (h *MskJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListClusters" || action == "AmazonMSK.ListClusters" {
		b := common.NewXmlBuilder()
		b.Start("ListClustersResponse").Start("ListClustersResult")
		b.Start("clusters").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
