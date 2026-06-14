package elasticache

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type ElastiCacheQueryHandler struct {
	service *ElastiCacheService
}

func NewElastiCacheQueryHandler(service *ElastiCacheService) *ElastiCacheQueryHandler {
	return &ElastiCacheQueryHandler{
		service: service,
	}
}

func (h *ElastiCacheQueryHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	switch action {
	case "DescribeCacheClusters":
		return map[string]any{"CacheClusters": []any{}}, nil
	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *ElastiCacheQueryHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "DescribeCacheClusters" {
		b := common.NewXmlBuilder()
		// Namespace?
		b.Start("DescribeCacheClustersResponse").Start("DescribeCacheClustersResult")
		b.Start("CacheClusters").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
