package stepfunctions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type StepFunctionsJsonHandler struct {
	service *StepFunctionsService
}

func NewStepFunctionsJsonHandler(service *StepFunctionsService) *StepFunctionsJsonHandler {
	return &StepFunctionsJsonHandler{
		service: service,
	}
}

func (h *StepFunctionsJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateStateMachine":
		var req struct {
			Name       string `json:"name"`
			Definition string `json:"definition"`
			RoleArn    string `json:"roleArn"`
		}
		json.Unmarshal(request, &req)
		sm, err := h.service.CreateStateMachine(ctx, req.Name, req.Definition, req.RoleArn)
		if err != nil {
			return nil, err
		}
		return map[string]string{"stateMachineArn": sm.StateMachineArn, "creationDate": fmt.Sprintf("%f", float64(sm.CreationDate.Unix()))}, nil

	case "DescribeStateMachine":
		var req struct {
			StateMachineArn string `json:"stateMachineArn"`
		}
		json.Unmarshal(request, &req)
		sm, err := h.service.DescribeStateMachine(ctx, req.StateMachineArn)
		if err != nil {
			return nil, err
		}
		return sm, nil

	case "ListStateMachines":
		sms, _ := h.service.ListStateMachines(ctx)
		return map[string]any{"stateMachines": sms}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *StepFunctionsJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("StepFunctions does not support Query protocol")
}
