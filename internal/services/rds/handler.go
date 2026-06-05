package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type RdsQueryHandler struct {
	service *RdsService
}

func NewRdsQueryHandler(service *RdsService) *RdsQueryHandler {
	return &RdsQueryHandler{
		service: service,
	}
}

func (h *RdsQueryHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("RDS does not support JSON protocol")
}

func (h *RdsQueryHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	ctx := context.Background()

	switch action {
	case "CreateDBInstance":
		id := params.Get("DBInstanceIdentifier")
		engine := params.Get("Engine")
		class := params.Get("DBInstanceClass")
		size, _ := strconv.Atoi(params.Get("AllocatedStorage"))
		inst, err := h.service.CreateDBInstance(ctx, id, engine, class, size)
		if err != nil {
			return "", err
		}
		
		b := common.NewXmlBuilder().Start("CreateDBInstanceResponse").Start("CreateDBInstanceResult")
		b.Start("DBInstance").
			Elem("DBInstanceIdentifier", inst.DBInstanceIdentifier).
			Elem("DBInstanceStatus", inst.DBInstanceStatus).
			End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-rds").End().End()
		return b.Build(), nil

	case "DescribeDBInstances":
		id := params.Get("DBInstanceIdentifier")
		instances, _ := h.service.DescribeDBInstances(ctx, id)
		
		b := common.NewXmlBuilder().Start("DescribeDBInstancesResponse").Start("DescribeDBInstancesResult").Start("DBInstances")
		for _, inst := range instances {
			b.Start("DBInstance").
				Elem("DBInstanceIdentifier", inst.DBInstanceIdentifier).
				Elem("DBInstanceStatus", inst.DBInstanceStatus).
				Elem("Engine", inst.Engine).
				Start("Endpoint").
				Elem("Address", inst.Endpoint.Address).
				Elem("Port", fmt.Sprintf("%d", inst.Endpoint.Port)).
				End().
				End()
		}
		b.End().End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-rds").End().End()
		return b.Build(), nil

	default:
		return "", fmt.Errorf("Unknown RDS action: %s", action)
	}
}
