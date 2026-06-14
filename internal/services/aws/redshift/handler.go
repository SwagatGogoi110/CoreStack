package redshift

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type RedshiftQueryHandler struct {
	service *RedshiftService
}

func NewRedshiftQueryHandler(service *RedshiftService) *RedshiftQueryHandler {
	return &RedshiftQueryHandler{
		service: service,
	}
}

func (h *RedshiftQueryHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("Redshift does not support JSON protocol")
}

func (h *RedshiftQueryHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	ctx := context.Background()

	switch action {
	case "CreateCluster":
		id := params.Get("ClusterIdentifier")
		nodeType := params.Get("NodeType")
		masterUser := params.Get("MasterUsername")
		dbName := params.Get("DBName")
		cluster, err := h.service.CreateCluster(ctx, id, nodeType, masterUser, dbName)
		if err != nil {
			return "", err
		}

		b := common.NewXmlBuilder().Start("CreateClusterResponse").Start("CreateClusterResult")
		b.Start("Cluster").
			Elem("ClusterIdentifier", cluster.ClusterIdentifier).
			Elem("ClusterStatus", cluster.ClusterStatus).
			Elem("NodeType", cluster.NodeType).
			End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-redshift").End().End()
		return b.Build(), nil

	case "DescribeClusters":
		id := params.Get("ClusterIdentifier")
		clusters, err := h.service.DescribeClusters(ctx, id)
		if err != nil {
			return "", err
		}

		b := common.NewXmlBuilder().Start("DescribeClustersResponse").Start("DescribeClustersResult").Start("Clusters")
		for _, c := range clusters {
			b.Start("Cluster").
				Elem("ClusterIdentifier", c.ClusterIdentifier).
				Elem("ClusterStatus", c.ClusterStatus).
				Elem("NodeType", c.NodeType).
				Elem("MasterUsername", c.MasterUsername).
				Elem("DBName", c.DBName).
				Start("Endpoint").
				Elem("Address", c.Endpoint.Address).
				Elem("Port", fmt.Sprintf("%d", c.Endpoint.Port)).
				End().
				End()
		}
		b.End().End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-redshift").End().End()
		return b.Build(), nil

	case "DeleteCluster":
		id := params.Get("ClusterIdentifier")
		err := h.service.DeleteCluster(ctx, id)
		if err != nil {
			return "", err
		}

		b := common.NewXmlBuilder().Start("DeleteClusterResponse").Start("DeleteClusterResult")
		b.Start("Cluster").
			Elem("ClusterIdentifier", id).
			Elem("ClusterStatus", "deleting").
			End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-redshift").End().End()
		return b.Build(), nil

	default:
		return "", fmt.Errorf("Unknown Redshift action: %s", action)
	}
}
