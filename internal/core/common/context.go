package common

import (
	"context"
)

type contextKey string

// RequestContextKey is the key for the request context.
const RequestContextKey contextKey = "CloudStack-request-context"

// RequestContext holds per-request derived values.
type RequestContext struct {
	AccountID string
	Region    string
}

// FromContext retrieves the RequestContext from the request context.
func FromContext(ctx context.Context) *RequestContext {
	val, ok := ctx.Value(RequestContextKey).(*RequestContext)
	if !ok {
		return &RequestContext{}
	}
	return val
}

// WithContext returns a new context with the RequestContext attached.
func WithContext(ctx context.Context, rc *RequestContext) context.Context {
	return context.WithValue(ctx, RequestContextKey, rc)
}

