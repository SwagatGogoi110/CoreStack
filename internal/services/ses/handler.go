package ses

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type SesQueryHandler struct {
	service *SesService
}

func NewSesQueryHandler(service *SesService) *SesQueryHandler {
	return &SesQueryHandler{
		service: service,
	}
}

func (h *SesQueryHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	switch action {
	case "ListIdentities":
		return map[string]any{"Identities": []string{}}, nil
	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *SesQueryHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	if action == "ListIdentities" {
		b := common.NewXmlBuilder()
		// Namespace?
		b.Start("ListIdentitiesResponse").Start("ListIdentitiesResult")
		b.Start("Identities").End()
		b.End().Start("ResponseMetadata").Elem("RequestId", "CloudStack").End().End()
		return b.Build(), nil
	}
	return "", fmt.Errorf("Unknown action: %s", action)
}
