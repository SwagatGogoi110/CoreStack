package acm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type AcmJsonHandler struct {
	service *AcmService
}

func NewAcmJsonHandler(service *AcmService) *AcmJsonHandler {
	return &AcmJsonHandler{
		service: service,
	}
}

func (h *AcmJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "RequestCertificate":
		var req struct {
			DomainName              string   `json:"DomainName"`
			SubjectAlternativeNames []string `json:"SubjectAlternativeNames"`
		}
		json.Unmarshal(request, &req)
		cert, err := h.service.RequestCertificate(ctx, req.DomainName, req.SubjectAlternativeNames)
		if err != nil {
			return nil, err
		}
		return map[string]string{"CertificateArn": cert.Arn}, nil

	case "DescribeCertificate":
		var req struct {
			CertificateArn string `json:"CertificateArn"`
		}
		json.Unmarshal(request, &req)
		cert, err := h.service.DescribeCertificate(ctx, req.CertificateArn)
		if err != nil {
			return nil, err
		}
		return map[string]any{"Certificate": cert}, nil

	case "ListCertificates":
		summaries, _ := h.service.ListCertificates(ctx)
		return map[string]any{"CertificateSummaryList": summaries}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *AcmJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("ACM does not support Query protocol")
}
