package docdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type DocdbQueryHandler struct {
	service *DocdbService
}

func NewDocdbQueryHandler(service *DocdbService) *DocdbQueryHandler {
	return &DocdbQueryHandler{
		service: service,
	}
}

func (h *DocdbQueryHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("DocumentDB does not support JSON protocol")
}

func (h *DocdbQueryHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	ctx := context.Background()

	switch action {
	case "CreateDBCluster":
		id := params.Get("DBClusterIdentifier")
		engine := params.Get("Engine")
		cluster, err := h.service.CreateDBCluster(ctx, id, engine)
		if err != nil {
			return "", err
		}

		b := common.NewXmlBuilder().Start("CreateDBClusterResponse").Start("CreateDBClusterResult")
		b.Start("DBCluster").
			Elem("DBClusterIdentifier", cluster.DBClusterIdentifier).
			Elem("Status", cluster.Status).
			Elem("Engine", cluster.Engine).
			End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-docdb").End().End()
		return b.Build(), nil

	case "DescribeDBClusters":
		id := params.Get("DBClusterIdentifier")
		clusters, err := h.service.DescribeDBClusters(ctx, id)
		if err != nil {
			return "", err
		}

		b := common.NewXmlBuilder().Start("DescribeDBClustersResponse").Start("DescribeDBClustersResult").Start("DBClusters")
		for _, c := range clusters {
			b.Start("DBCluster").
				Elem("DBClusterIdentifier", c.DBClusterIdentifier).
				Elem("Status", c.Status).
				Elem("Engine", c.Engine).
				Elem("Endpoint", c.Endpoint).
				Elem("Port", fmt.Sprintf("%d", c.Port)).
				End()
		}
		b.End().End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-docdb").End().End()
		return b.Build(), nil

	case "DeleteDBCluster":
		id := params.Get("DBClusterIdentifier")
		err := h.service.DeleteDBCluster(ctx, id)
		if err != nil {
			return "", err
		}

		b := common.NewXmlBuilder().Start("DeleteDBClusterResponse").Start("DeleteDBClusterResult")
		b.Start("DBCluster").
			Elem("DBClusterIdentifier", id).
			Elem("Status", "deleting").
			End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-docdb").End().End()
		return b.Build(), nil

	default:
		return "", fmt.Errorf("Unknown DocumentDB action: %s", action)
	}
}
