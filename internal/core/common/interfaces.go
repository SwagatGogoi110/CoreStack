package common

import (
	"encoding/json"
	"net/url"
)

// ServiceHandler defines the interface for an AWS service handler.
type ServiceHandler interface {
	HandleJSON(action string, request json.RawMessage, rc *RequestContext) (any, error)
	HandleQuery(action string, params url.Values, rc *RequestContext) (string, error)
}
