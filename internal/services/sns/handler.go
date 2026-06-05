package sns

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type SnsQueryHandler struct {
	service *SnsService
}

func NewSnsQueryHandler(service *SnsService) *SnsQueryHandler {
	return &SnsQueryHandler{
		service: service,
	}
}

func (h *SnsQueryHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("SNS JSON protocol not implemented in this handler")
}

func (h *SnsQueryHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	ctx := context.Background()

	switch action {
	case "CreateTopic":
		name := params.Get("Name")
		t, err := h.service.CreateTopic(ctx, name, nil)
		if err != nil {
			return "", err
		}
		return h.xmlResponse("CreateTopic", map[string]string{"TopicArn": t.TopicArn}), nil

	case "ListTopics":
		topics, _ := h.service.ListTopics(ctx)
		b := common.NewXmlBuilder().Start("ListTopicsResponse").Start("ListTopicsResult").Start("Topics")
		for _, t := range topics {
			b.Start("member").Elem("TopicArn", t.TopicArn).End()
		}
		b.End().End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-sns").End().End()
		return b.Build(), nil

	case "Subscribe":
		topicArn := params.Get("TopicArn")
		protocol := params.Get("Protocol")
		endpoint := params.Get("Endpoint")
		sub, err := h.service.Subscribe(ctx, topicArn, protocol, endpoint)
		if err != nil {
			return "", err
		}
		return h.xmlResponse("Subscribe", map[string]string{"SubscriptionArn": sub.SubscriptionArn}), nil

	case "Publish":
		topicArn := params.Get("TopicArn")
		message := params.Get("Message")
		id, err := h.service.Publish(ctx, topicArn, message)
		if err != nil {
			return "", err
		}
		return h.xmlResponse("Publish", map[string]string{"MessageId": id}), nil

	default:
		return "", fmt.Errorf("Unknown SNS action: %s", action)
	}
}

func (h *SnsQueryHandler) xmlResponse(action string, fields map[string]string) string {
	b := common.NewXmlBuilder().Start(action + "Response").Start(action + "Result")
	for k, v := range fields {
		b.Elem(k, v)
	}
	b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-sns").End().End()
	return b.Build()
}
