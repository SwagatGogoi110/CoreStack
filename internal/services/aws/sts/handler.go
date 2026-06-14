package sts

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/core/common"
	"github.com/hectorvent/cloudstack/internal/services/aws/iam"
)

type StsQueryHandler struct {
	iamService *iam.IamService
}

func NewStsQueryHandler(iamService *iam.IamService) *StsQueryHandler {
	return &StsQueryHandler{
		iamService: iamService,
	}
}

func (h *StsQueryHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("STS does not support JSON protocol")
}

func (h *StsQueryHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	switch action {
	case "GetCallerIdentity":
		return h.handleGetCallerIdentity(rc), nil
	case "GetSessionToken":
		return h.handleGetSessionToken(params), nil
	default:
		return "", fmt.Errorf("Unknown STS action: %s", action)
	}
}

func (h *StsQueryHandler) handleGetCallerIdentity(rc *common.RequestContext) string {
	accountID := rc.AccountID
	if accountID == "" {
		accountID = "000000000000"
	}

	result := common.NewXmlBuilder().
		Elem("UserId", accountID).
		Elem("Account", accountID).
		Elem("Arn", fmt.Sprintf("arn:aws:iam::%s:root", accountID)).
		Build()

	return h.envelope("GetCallerIdentity", result)
}

func (h *StsQueryHandler) handleGetSessionToken(params url.Values) string {
	accessKeyID := "ASIA" + h.randomID(16)
	secretKey := h.randomSecret(40)
	sessionToken := h.randomSecret(200)
	expiration := time.Now().Add(12 * time.Hour)

	result := h.credentialsXml(accessKeyID, secretKey, sessionToken, expiration)
	return h.envelope("GetSessionToken", result)
}

func (h *StsQueryHandler) credentialsXml(accessKeyID, secretKey, sessionToken string, expiration time.Time) string {
	return common.NewXmlBuilder().
		Start("Credentials").
		Elem("AccessKeyId", accessKeyID).
		Elem("SecretAccessKey", secretKey).
		Elem("SessionToken", sessionToken).
		Elem("Expiration", expiration.Format(time.RFC3339)).
		End().
		Build()
}

func (h *StsQueryHandler) envelope(action, result string) string {
	return common.NewXmlBuilder().
		Start(action + "Response").
		Start(action + "Result").
		Raw(result).
		End().
		Start("ResponseMetadata").
		Elem("RequestId", uuid.New().String()).
		End().
		End().
		Build()
}

func (h *StsQueryHandler) randomID(length int) string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func (h *StsQueryHandler) randomSecret(length int) string {
	secretChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	b := make([]byte, length)
	for i := range b {
		b[i] = secretChars[rand.Intn(len(secretChars))]
	}
	return string(b)
}
