package cognito

import (
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
	// Handle with and without prefix
	if action == "ListUserPools" || action == "AWSCognitoIdentityProviderService.ListUserPools" {
		return map[string]any{"UserPools": []any{}}, nil
	}
	return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
}

func (h *CognitoJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListUserPools" || action == "AWSCognitoIdentityProviderService.ListUserPools" {
		b := common.NewXmlBuilder()
		b.Start("ListUserPoolsResponse").Start("ListUserPoolsResult")
		b.Start("UserPools").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
