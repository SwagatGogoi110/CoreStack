package scheduler

import (
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
	// Handle with and without prefix
	if action == "ListSchedules" || action == "AmazonScheduler.ListSchedules" {
		return map[string]any{"Schedules": []any{}}, nil
	}
	return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
}

func (h *SchedulerJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListSchedules" || action == "AmazonScheduler.ListSchedules" {
		b := common.NewXmlBuilder()
		b.Start("ListSchedulesResponse").Start("ListSchedulesResult")
		b.Start("Schedules").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
