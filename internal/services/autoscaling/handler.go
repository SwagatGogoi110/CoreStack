package autoscaling

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
	"github.com/hectorvent/cloudstack/internal/services/autoscaling/model"
)

type AutoScalingQueryHandler struct {
	service *AutoScalingService
}

func NewAutoScalingQueryHandler(service *AutoScalingService) *AutoScalingQueryHandler {
	return &AutoScalingQueryHandler{
		service: service,
	}
}

func (h *AutoScalingQueryHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("AutoScaling does not support JSON protocol")
}

func (h *AutoScalingQueryHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	ctx := context.Background()

	switch action {
	case "CreateAutoScalingGroup":
		asg := &model.AutoScalingGroup{
			AutoScalingGroupName:    params.Get("AutoScalingGroupName"),
			LaunchConfigurationName: params.Get("LaunchConfigurationName"),
			// ... other params
		}
		if err := h.service.CreateAutoScalingGroup(ctx, asg); err != nil {
			return "", err
		}
		return h.xmlResponse("CreateAutoScalingGroup", nil), nil

	case "DescribeAutoScalingGroups":
		groups, _ := h.service.DescribeAutoScalingGroups(ctx, nil)
		b := common.NewXmlBuilder().Start("DescribeAutoScalingGroupsResponse").Start("DescribeAutoScalingGroupsResult").Start("AutoScalingGroups")
		for _, g := range groups {
			b.Start("member").
				Elem("AutoScalingGroupName", g.AutoScalingGroupName).
				Elem("AutoScalingGroupARN", g.AutoScalingGroupArn).
				End()
		}
		b.End().End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-asg").End().End()
		return b.Build(), nil

	case "CreateLaunchConfiguration":
		lc := &model.LaunchConfiguration{
			LaunchConfigurationName: params.Get("LaunchConfigurationName"),
			ImageID:                 params.Get("ImageId"),
			InstanceType:            params.Get("InstanceType"),
		}
		if err := h.service.CreateLaunchConfiguration(ctx, lc); err != nil {
			return "", err
		}
		return h.xmlResponse("CreateLaunchConfiguration", nil), nil

	default:
		return "", fmt.Errorf("Unknown AutoScaling action: %s", action)
	}
}

func (h *AutoScalingQueryHandler) xmlResponse(action string, fields map[string]string) string {
	b := common.NewXmlBuilder().Start(action + "Response")
	if len(fields) > 0 {
		b.Start(action + "Result")
		for k, v := range fields {
			b.Elem(k, v)
		}
		b.End()
	}
	b.Start("ResponseMetadata").Elem("RequestId", "CloudStack-asg").End().End()
	return b.Build()
}
