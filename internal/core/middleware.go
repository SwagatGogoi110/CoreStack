package core

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/core/common"
)

// RequestContextMiddleware populates the RequestContext in the request.
func (s *Server) RequestContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		region := s.resolver.ResolveRegion(auth)
		accountID := s.resolver.ResolveAccount(auth)

		rc := &common.RequestContext{
			AccountID: accountID,
			Region:    region,
		}

		ctx := common.WithContext(r.Context(), rc)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


const (
	amzRequestID  = "x-amz-request-id"
	amznRequestID = "x-amzn-RequestId"
	amzID2        = "x-amz-id-2"
)

// RequestIDMiddleware adds AWS request-id response headers to every HTTP response.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()

		// Intercept the response headers to add the request ID if not present.
		// In Go, we usually set headers before calling next.ServeHTTP if they are fixed,
		// but since these can be overridden by handlers, we should wrap the ResponseWriter
		// or just set them if they aren't there after.
		// However, AWS SDKs expect these headers.
		
		w.Header().Set(amzRequestID, requestID)
		w.Header().Set(amznRequestID, requestID)
		w.Header().Set(amzID2, requestID)

		next.ServeHTTP(w, r)
	})
}

// GlobalCorsFilterMiddleware handles CORS headers.
func GlobalCorsFilterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Expose-Headers", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
