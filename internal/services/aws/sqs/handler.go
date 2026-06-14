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
	ctx := context.Background()

	switch action {
	case "CreateQueue":
		var req struct {
			QueueName  string            `json:"QueueName"`
			Attributes map[string]string `json:"Attributes"`
		}
		json.Unmarshal(request, &req)
		q, err := h.service.CreateQueue(ctx, req.QueueName, req.Attributes)
		if err != nil {
			return nil, err
		}
		return map[string]string{"QueueUrl": q.QueueURL}, nil

	case "GetQueueUrl":
		var req struct {
			QueueName string `json:"QueueName"`
		}
		json.Unmarshal(request, &req)
		url, err := h.service.GetQueueUrl(ctx, req.QueueName)
		if err != nil {
			return nil, err
		}
		return map[string]string{"QueueUrl": url}, nil

	case "ListQueues":
		var req struct {
			QueueNamePrefix string `json:"QueueNamePrefix"`
		}
		json.Unmarshal(request, &req)
		urls, _ := h.service.ListQueues(ctx, req.QueueNamePrefix)
		return map[string]any{"QueueUrls": urls}, nil

	case "SendMessage":
		var req struct {
			QueueUrl    string `json:"QueueUrl"`
			MessageBody string `json:"MessageBody"`
		}
		json.Unmarshal(request, &req)
		queueName := h.extractQueueName(req.QueueUrl)
		msg, err := h.service.SendMessage(ctx, queueName, req.MessageBody)
		if err != nil {
			return nil, err
		}
		return map[string]string{
			"MessageId":        msg.MessageID,
			"MD5OfMessageBody": msg.MD5OfBody,
		}, nil

	case "ReceiveMessage":
		var req struct {
			QueueUrl            string `json:"QueueUrl"`
			MaxNumberOfMessages int    `json:"MaxNumberOfMessages"`
		}
		json.Unmarshal(request, &req)
		queueName := h.extractQueueName(req.QueueUrl)
		msgs, err := h.service.ReceiveMessage(ctx, queueName, req.MaxNumberOfMessages)
		if err != nil {
			return nil, err
		}
		return map[string]any{"Messages": msgs}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
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

	case "ListQueues":
		prefix := params.Get("QueueNamePrefix")
		urls, _ := h.service.ListQueues(ctx, prefix)
		
		b := common.NewXmlBuilder()
		b.Raw("<ListQueuesResponse xmlns=\"http://queue.amazonaws.com/doc/2012-11-05/\">")
		b.Start("ListQueuesResult")
		for _, u := range urls {
			b.Elem("QueueUrl", u)
		}
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-sqs").End().Raw("</ListQueuesResponse>")
		return b.Build(), nil

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
	b := common.NewXmlBuilder()
	b.Raw(fmt.Sprintf("<%sResponse xmlns=\"http://queue.amazonaws.com/doc/2012-11-05/\">", action))
	b.Start(action + "Result")
	for k, v := range fields {
		b.Elem(k, v)
	}
	b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-sqs").End().Raw(fmt.Sprintf("</%sResponse>", action))
	return b.Build()
}
