package neptune

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type NeptuneQueryHandler struct {
	service *NeptuneService
}

func NewNeptuneQueryHandler(service *NeptuneService) *NeptuneQueryHandler {
	return &NeptuneQueryHandler{
		service: service,
	}
}

func (h *NeptuneQueryHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("Neptune does not support JSON protocol")
}

func (h *NeptuneQueryHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	ctx := context.Background()

	switch action {
	case "CreateDBCluster":
		id := params.Get("DBClusterIdentifier")
		version := params.Get("EngineVersion")
		c, err := h.service.CreateDBCluster(ctx, id, version)
		if err != nil {
			return "", err
		}
		
		b := common.NewXmlBuilder().Start("CreateDBClusterResponse").Start("CreateDBClusterResult")
		b.Start("DBCluster").
			Elem("DBClusterIdentifier", c.DBClusterIdentifier).
			Elem("Status", c.Status).
			End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-neptune").End().End()
		return b.Build(), nil

	case "DescribeDBClusters":
		id := params.Get("DBClusterIdentifier")
		clusters, _ := h.service.DescribeDBClusters(ctx, id)
		
		b := common.NewXmlBuilder().Start("DescribeDBClustersResponse").Start("DescribeDBClustersResult").Start("DBClusters")
		for _, c := range clusters {
			b.Start("DBCluster").
				Elem("DBClusterIdentifier", c.DBClusterIdentifier).
				Elem("Status", c.Status).
				Elem("Endpoint", c.Endpoint).
				Elem("Port", fmt.Sprintf("%d", c.Port)).
				End()
		}
		b.End().End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-neptune").End().End()
		return b.Build(), nil

	default:
		return "", fmt.Errorf("Unknown Neptune action: %s", action)
	}
}
