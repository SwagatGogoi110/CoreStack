package cognito

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type CognitoJsonHandler struct {
	service *CognitoService
}

func NewCognitoJsonHandler(service *CognitoService) *CognitoJsonHandler {
	return &CognitoJsonHandler{
		service: service,
	}
}

func (h *CognitoJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateUserPool":
		var req struct {
			PoolName string `json:"PoolName"`
		}
		if err := json.Unmarshal(request, &req); err != nil {
			return nil, err
		}
		pool, err := h.service.CreateUserPool(ctx, req.PoolName)
		if err != nil {
			return nil, err
		}
		return map[string]any{
			"UserPool": map[string]any{
				"Id":   pool.ID,
				"Name": pool.Name,
				"Arn":  pool.Arn,
			},
		}, nil

	case "DescribeUserPool":
		var req struct {
			UserPoolID string `json:"UserPoolId"`
		}
		if err := json.Unmarshal(request, &req); err != nil {
			return nil, err
		}
		pool, err := h.service.DescribeUserPool(ctx, req.UserPoolID)
		if err != nil {
			return nil, err
		}
		return map[string]any{
			"UserPool": map[string]any{
				"Id":   pool.ID,
				"Name": pool.Name,
				"Arn":  pool.Arn,
			},
		}, nil

	case "CreateUserPoolClient":
		var req struct {
			UserPoolID string `json:"UserPoolId"`
			ClientName string `json:"ClientName"`
		}
		if err := json.Unmarshal(request, &req); err != nil {
			return nil, err
		}
		client, err := h.service.CreateUserPoolClient(ctx, req.UserPoolID, req.ClientName)
		if err != nil {
			return nil, err
		}
		return map[string]any{
			"UserPoolClient": map[string]any{
				"ClientId":   client.ClientID,
				"UserPoolId": client.UserPoolID,
				"ClientName": client.ClientName,
			},
		}, nil

	case "AdminCreateUser":
		var req struct {
			UserPoolID     string            `json:"UserPoolId"`
			Username       string            `json:"Username"`
			UserAttributes []map[string]string `json:"UserAttributes"`
		}
		if err := json.Unmarshal(request, &req); err != nil {
			return nil, err
		}
		attrs := make(map[string]string)
		for _, a := range req.UserAttributes {
			attrs[a["Name"]] = a["Value"]
		}
		user, err := h.service.AdminCreateUser(ctx, req.UserPoolID, req.Username, attrs)
		if err != nil {
			return nil, err
		}
		
		resAttrs := make([]map[string]string, 0)
		for k, v := range user.Attributes {
			resAttrs = append(resAttrs, map[string]string{"Name": k, "Value": v})
		}

		return map[string]any{
			"User": map[string]any{
				"Username":       user.Username,
				"Attributes":     resAttrs,
				"UserCreateDate": user.CreationDate,
				"UserStatus":     user.UserStatus,
				"Enabled":        user.Enabled,
			},
		}, nil

	default:
		return nil, fmt.Errorf("Unknown Cognito action: %s", action)
	}
}

func (h *CognitoJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("Cognito does not support Query protocol")
}
