package opensearch

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type OpenSearchJsonHandler struct {
	service *OpenSearchService
}

func NewOpenSearchJsonHandler(service *OpenSearchService) *OpenSearchJsonHandler {
	return &OpenSearchJsonHandler{
		service: service,
	}
}

func (h *OpenSearchJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	// Handle with and without prefix
	if action == "ListDomainNames" || action == "AmazonOpenSearchService.ListDomainNames" {
		return map[string]any{"DomainNames": []any{}}, nil
	}
	return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
}

func (h *OpenSearchJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListDomainNames" || action == "AmazonOpenSearchService.ListDomainNames" {
		b := common.NewXmlBuilder()
		b.Start("ListDomainNamesResponse").Start("ListDomainNamesResult")
		b.Start("DomainNames").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
