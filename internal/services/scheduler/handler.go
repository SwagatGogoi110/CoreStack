package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type SchedulerJsonHandler struct {
	service *SchedulerService
}

func NewSchedulerJsonHandler(service *SchedulerService) *SchedulerJsonHandler {
	return &SchedulerJsonHandler{
		service: service,
	}
}

func (h *SchedulerJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateSchedule":
		var req struct {
			Name               string `json:"Name"`
			GroupName          string `json:"GroupName"`
			ScheduleExpression string `json:"ScheduleExpression"`
		}
		json.Unmarshal(request, &req)
		schedule, err := h.service.CreateSchedule(ctx, req.Name, req.GroupName, req.ScheduleExpression)
		if err != nil {
			return nil, err
		}
		return map[string]string{"ScheduleArn": schedule.Arn}, nil

	case "GetSchedule":
		var req struct {
			Name      string `json:"Name"`
			GroupName string `json:"GroupName"`
		}
		json.Unmarshal(request, &req)
		schedule, err := h.service.GetSchedule(ctx, req.Name, req.GroupName)
		if err != nil {
			return nil, err
		}
		return schedule, nil

	case "ListSchedules":
		schedules, _ := h.service.ListSchedules(ctx)
		return map[string]any{"Schedules": schedules}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *SchedulerJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("Scheduler does not support Query protocol")
}
