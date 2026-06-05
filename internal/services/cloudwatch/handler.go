package cloudwatch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/hectorvent/cloudstack/internal/core/common"
	"github.com/hectorvent/cloudstack/internal/services/cloudwatch/model"
)

type CloudWatchQueryHandler struct {
	service *CloudWatchService
}

func NewCloudWatchQueryHandler(service *CloudWatchService) *CloudWatchQueryHandler {
	return &CloudWatchQueryHandler{
		service: service,
	}
}

func (h *CloudWatchQueryHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("CloudWatch does not support JSON protocol")
}

func (h *CloudWatchQueryHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	ctx := context.Background()

	switch action {
	case "PutMetricData":
		namespace := params.Get("Namespace")
		// Query protocol for list of datums is complex (MetricData.member.1.MetricName ...)
		// For now, let's just log it and return OK.
		h.service.PutMetricData(ctx, namespace, nil)
		return h.xmlResponse("PutMetricData", nil), nil

	case "PutMetricAlarm":
		threshold, _ := strconv.ParseFloat(params.Get("Threshold"), 64)
		alarm := &model.MetricAlarm{
			AlarmName:          params.Get("AlarmName"),
			MetricName:         params.Get("MetricName"),
			Namespace:          params.Get("Namespace"),
			Threshold:          threshold,
			ComparisonOperator: params.Get("ComparisonOperator"),
		}
		if err := h.service.PutMetricAlarm(ctx, alarm); err != nil {
			return "", err
		}
		return h.xmlResponse("PutMetricAlarm", nil), nil

	case "DescribeAlarms":
		alarms, _ := h.service.DescribeAlarms(ctx)
		b := common.NewXmlBuilder().Start("DescribeAlarmsResponse").Start("DescribeAlarmsResult").Start("MetricAlarms")
		for _, a := range alarms {
			b.Start("member").
				Elem("AlarmName", a.AlarmName).
				Elem("AlarmArn", a.AlarmArn).
				Elem("StateValue", a.StateValue).
				End()
		}
		b.End().End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-cw").End().End()
		return b.Build(), nil

	default:
		return "", fmt.Errorf("Unknown CloudWatch action: %s", action)
	}
}

func (h *CloudWatchQueryHandler) xmlResponse(action string, fields map[string]string) string {
	b := common.NewXmlBuilder().Start(action + "Response")
	if len(fields) > 0 {
		b.Start(action + "Result")
		for k, v := range fields {
			b.Elem(k, v)
		}
		b.End()
	}
	b.Start("ResponseMetadata").Elem("RequestId", "CloudStack-cw").End().End()
	return b.Build()
}
