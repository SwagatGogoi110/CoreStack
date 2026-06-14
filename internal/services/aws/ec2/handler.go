package ec2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type Ec2QueryHandler struct {
	service *Ec2Service
}

func NewEc2QueryHandler(service *Ec2Service) *Ec2QueryHandler {
	return &Ec2QueryHandler{
		service: service,
	}
}

func (h *Ec2QueryHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("EC2 does not support JSON protocol")
}

func (h *Ec2QueryHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	ctx := context.Background()

	switch action {
	case "RunInstances":
		imageId := params.Get("ImageId")
		instType := params.Get("InstanceType")
		min, _ := strconv.Atoi(params.Get("MinCount"))
		if min == 0 { min = 1 }
		res, err := h.service.RunInstances(ctx, imageId, instType, min)
		if err != nil {
			return "", err
		}
		
		b := common.NewXmlBuilder().Start("RunInstancesResponse")
		b.Elem("requestId", "CloudStack-ec2")
		b.Elem("reservationId", res.ReservationID)
		b.Elem("ownerId", res.OwnerID)
		b.Start("instancesSet")
		for _, i := range res.Instances {
			b.Start("item").
				Elem("instanceId", i.InstanceID).
				Elem("imageId", i.ImageID).
				Elem("instanceType", i.InstanceType).
				Start("instanceState").Elem("code", "16").Elem("name", i.State).End().
				Elem("privateIpAddress", i.PrivateIp).
				End()
		}
		b.End().End()
		return b.Build(), nil

	case "DescribeInstances":
		reservations, _ := h.service.DescribeInstances(ctx)
		b := common.NewXmlBuilder().Start("DescribeInstancesResponse")
		b.Elem("requestId", "CloudStack-ec2")
		b.Start("reservationSet")
		for _, r := range reservations {
			b.Start("item").
				Elem("reservationId", r.ReservationID).
				Elem("ownerId", r.OwnerID).
				Start("instancesSet")
			for _, i := range r.Instances {
				b.Start("item").
					Elem("instanceId", i.InstanceID).
					Elem("imageId", i.ImageID).
					Elem("instanceType", i.InstanceType).
					Start("instanceState").Elem("code", "16").Elem("name", i.State).End().
					End()
			}
			b.End().End()
		}
		b.End().End()
		return b.Build(), nil

	default:
		return "", fmt.Errorf("Unknown EC2 action: %s", action)
	}
}
