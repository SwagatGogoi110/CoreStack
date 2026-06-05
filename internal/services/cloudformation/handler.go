package cloudformation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type CloudFormationQueryHandler struct {
	service *CloudFormationService
}

func NewCloudFormationQueryHandler(service *CloudFormationService) *CloudFormationQueryHandler {
	return &CloudFormationQueryHandler{
		service: service,
	}
}

func (h *CloudFormationQueryHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("CloudFormation does not support JSON protocol")
}

func (h *CloudFormationQueryHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	ctx := context.Background()

	switch action {
	case "CreateStack":
		name := params.Get("StackName")
		template := params.Get("TemplateBody")
		stack, err := h.service.CreateStack(ctx, name, template)
		if err != nil {
			return "", err
		}
		return h.xmlResponse("CreateStack", map[string]string{"StackId": stack.StackID}), nil

	case "DescribeStacks":
		name := params.Get("StackName")
		stacks, err := h.service.DescribeStacks(ctx, name)
		if err != nil {
			return "", err
		}
		
		b := common.NewXmlBuilder().Start("DescribeStacksResponse").Start("DescribeStacksResult").Start("Stacks")
		for _, s := range stacks {
			b.Start("member").
				Elem("StackId", s.StackID).
				Elem("StackName", s.StackName).
				Elem("StackStatus", s.StackStatus).
				Elem("CreationTime", s.CreationTime.Format("2006-01-02T15:04:05.000Z")).
				End()
		}
		b.End().End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-cf").End().End()
		return b.Build(), nil

	case "DeleteStack":
		name := params.Get("StackName")
		if err := h.service.DeleteStack(ctx, name); err != nil {
			return "", err
		}
		return h.xmlResponse("DeleteStack", nil), nil

	default:
		return "", fmt.Errorf("Unknown CloudFormation action: %s", action)
	}
}

func (h *CloudFormationQueryHandler) xmlResponse(action string, fields map[string]string) string {
	b := common.NewXmlBuilder().Start(action + "Response")
	if len(fields) > 0 {
		b.Start(action + "Result")
		for k, v := range fields {
			b.Elem(k, v)
		}
		b.End()
	}
	b.Start("ResponseMetadata").Elem("RequestId", "CloudStack-cf").End().End()
	return b.Build()
}
