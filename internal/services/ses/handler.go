package ses

import (
	"context"
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
	return nil, fmt.Errorf("SES does not support JSON protocol")
}

func (h *SesQueryHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	ctx := context.Background()

	switch action {
	case "VerifyEmailIdentity":
		email := params.Get("EmailAddress")
		_, err := h.service.VerifyEmailIdentity(ctx, email)
		if err != nil {
			return "", err
		}
		return h.xmlResponse("VerifyEmailIdentity", nil), nil

	case "ListIdentities":
		identities, _ := h.service.ListIdentities(ctx)
		b := common.NewXmlBuilder().Start("ListIdentitiesResponse").Start("ListIdentitiesResult").Start("Identities")
		for _, i := range identities {
			b.Elem("member", i)
		}
		b.End().End().Start("ResponseMetadata").Elem("RequestId", "CloudStack-ses").End().End()
		return b.Build(), nil

	case "SendEmail":
		source := params.Get("Source")
		subject := params.Get("Message.Subject.Data")
		// ... more complex param parsing for destinations
		id, err := h.service.SendEmail(ctx, source, []string{}, subject)
		if err != nil {
			return "", err
		}
		return h.xmlResponse("SendEmail", map[string]string{"MessageId": id}), nil

	default:
		return "", fmt.Errorf("Unknown SES action: %s", action)
	}
}

func (h *SesQueryHandler) xmlResponse(action string, fields map[string]string) string {
	b := common.NewXmlBuilder().Start(action + "Response")
	if len(fields) > 0 {
		b.Start(action + "Result")
		for k, v := range fields {
			b.Elem(k, v)
		}
		b.End()
	}
	b.Start("ResponseMetadata").Elem("RequestId", "CloudStack-ses").End().End()
	return b.Build()
}
