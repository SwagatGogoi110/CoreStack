package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
	"github.com/hectorvent/cloudstack/internal/services/aws/iam/model"
)

type IamQueryHandler struct {
	service *IamService
}

func NewIamQueryHandler(service *IamService) *IamQueryHandler {
	return &IamQueryHandler{
		service: service,
	}
}

func (h *IamQueryHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("IAM does not support JSON protocol")
}

func (h *IamQueryHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	ctx := context.Background() // TODO: use rc to populate context if needed
	
	switch action {
	case "CreateUser":
		userName := params.Get("UserName")
		path := params.Get("Path")
		user, err := h.service.CreateUser(ctx, userName, path)
		if err != nil {
			return "", err
		}
		return h.xmlResponse("CreateUser", user), nil
	case "GetUser":
		userName := params.Get("UserName")
		user, err := h.service.GetUser(ctx, userName)
		if err != nil {
			return "", err
		}
		return h.xmlResponse("GetUser", user), nil
	case "ListUsers":
		users, err := h.service.ListUsers(ctx)
		if err != nil {
			return "", err
		}
		return h.xmlResponse("ListUsers", users), nil
	case "DeleteUser":
		userName := params.Get("UserName")
		if err := h.service.DeleteUser(ctx, userName); err != nil {
			return "", err
		}
		return h.xmlResponse("DeleteUser", nil), nil
	case "CreateAccessKey":
		userName := params.Get("UserName")
		key, err := h.service.CreateAccessKey(ctx, userName)
		if err != nil {
			return "", err
		}
		return h.xmlResponse("CreateAccessKey", key), nil
	default:
		return "", fmt.Errorf("Unknown IAM action: %s", action)
	}
}

func (h *IamQueryHandler) xmlResponse(action string, data any) string {
	b := common.NewXmlBuilder()
	b.Raw(fmt.Sprintf("<%sResponse xmlns=\"https://iam.amazonaws.com/doc/2010-05-08/\">", action))
	b.Start(action + "Result")

	switch v := data.(type) {
	case *model.User:
		b.Start("User").
			Elem("UserId", v.UserID).
			Elem("UserName", v.UserName).
			Elem("Arn", v.Arn).
			Elem("Path", v.Path).
			Elem("CreateDate", v.CreateDate.Format("2006-01-02T15:04:05Z")).
			End()
	case []*model.User:
		b.Start("Users")
		for _, user := range v {
			b.Start("member").
				Elem("UserId", user.UserID).
				Elem("UserName", user.UserName).
				Elem("Arn", user.Arn).
				Elem("Path", user.Path).
				Elem("CreateDate", user.CreateDate.Format("2006-01-02T15:04:05Z")).
				End()
		}
		b.End().Elem("IsTruncated", "false")
	case *model.AccessKey:
		b.Start("AccessKey").
			Elem("AccessKeyId", v.AccessKeyID).
			Elem("SecretAccessKey", v.SecretAccessKey).
			Elem("Status", v.Status).
			Elem("UserName", v.UserName).
			Elem("CreateDate", v.CreateDate.Format("2006-01-02T15:04:05Z")).
			End()
	}

	b.End(). // Result
		Start("ResponseMetadata").
		Elem("RequestId", "CloudStack-request-id").
		End().
		Raw(fmt.Sprintf("</%sResponse>", action))

	return b.Build()
}
