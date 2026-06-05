package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type SqsQueryHandler struct {
	service *SqsService
}

func NewSqsQueryHandler(service *SqsService) *SqsQueryHandler {
	return &SqsQueryHandler{
		service: service,
	}
}

func (h *SqsQueryHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("SQS JSON protocol not implemented in this handler")
}

func (h *SqsQueryHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	ctx := context.Background()
	
	switch action {
	case "CreateQueue":
		name := params.Get("QueueName")
		q, err := h.service.CreateQueue(ctx, name, nil)
		if err != nil {
			return "", err
		}
		return h.xmlResponse("CreateQueue", map[string]string{"QueueUrl": q.QueueURL}), nil

	case "GetQueueUrl":
		name := params.Get("QueueName")
		url, err := h.service.GetQueueUrl(ctx, name)
		if err != nil {
			return "", err
		}
		return h.xmlResponse("GetQueueUrl", map[string]string{"QueueUrl": url}), nil

	case "SendMessage":
		queueUrl := params.Get("QueueUrl")
		queueName := h.extractQueueName(queueUrl)
		body := params.Get("MessageBody")
		msg, err := h.service.SendMessage(ctx, queueName, body)
		if err != nil {
			return "", err
		}
		return h.xmlResponse("SendMessage", map[string]string{
			"MessageId": msg.MessageID,
			"MD5OfMessageBody": msg.MD5OfBody,
		}), nil

	case "ReceiveMessage":
		queueUrl := params.Get("QueueUrl")
		queueName := h.extractQueueName(queueUrl)
		max, _ := strconv.Atoi(params.Get("MaxNumberOfMessages"))
		if max == 0 {
			max = 1
		}
		msgs, err := h.service.ReceiveMessage(ctx, queueName, max)
		if err != nil {
			return "", err
		}
		
		b := common.NewXmlBuilder().Start("ReceiveMessageResponse").Start("ReceiveMessageResult")
		for _, m := range msgs {
			b.Start("Message").
				Elem("MessageId", m.MessageID).
				Elem("ReceiptHandle", m.ReceiptHandle).
				Elem("MD5OfBody", m.MD5OfBody).
				Elem("Body", m.Body).
				End()
		}
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-sqs").End().End()
		return b.Build(), nil

	default:
		return "", fmt.Errorf("Unknown SQS action: %s", action)
	}
}

func (h *SqsQueryHandler) extractQueueName(queueUrl string) string {
	parts := strings.Split(queueUrl, "/")
	return parts[len(parts)-1]
}

func (h *SqsQueryHandler) xmlResponse(action string, fields map[string]string) string {
	b := common.NewXmlBuilder().Start(action + "Response").Start(action + "Result")
	for k, v := range fields {
		b.Elem(k, v)
	}
	b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-sqs").End().End()
	return b.Build()
}
