package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/core/common"
)

var servicePattern = regexp.MustCompile(`Credential=\S+/\d{8}/[^/]+/([^/]+)/`)

func (s *Server) handleAWSRequest(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	// 1. Check for Lambda REST API
	if strings.HasPrefix(r.URL.Path, "/2015-03-31") {
		s.lambda.ServeHTTP(w, r)
		return
	}

	// 2. Check for API Gateway REST API
	if strings.HasPrefix(r.URL.Path, "/restapis") {
		s.apigateway.ServeHTTP(w, r)
		return
	}

	// 3. Check for Bedrock Runtime
	if strings.HasPrefix(r.URL.Path, "/model/") {
		s.bedrockruntime.ServeHTTP(w, r)
		return
	}

	// 4. Check for CloudFront
	if strings.HasPrefix(r.URL.Path, "/2020-05-31") {
		s.cloudfront.ServeHTTP(w, r)
		return
	}

	// 5. Check for Route53
	if strings.HasPrefix(r.URL.Path, "/2013-04-01") {
		s.route53.ServeHTTP(w, r)
		return
	}

	// 6. Check for JSON protocols
	if contentType == "application/x-amz-json-1.0" || contentType == "application/x-amz-json-1.1" {
		s.handleJSONRequest(w, r)
		return
	}


	// 2. Check for Query protocol
	if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		// Distinguish between SQS/SNS and S3/others
		// SQS/SNS usually have an Action parameter
		if r.FormValue("Action") != "" {
			s.handleQueryRequest(w, r)
			return
		}
	}

	// 3. Fallback to S3 for REST-style or missing Content-Type
	// This is where S3 usually lands.
	s.s3.ServeHTTP(w, r)
}



func (s *Server) handleQueryRequest(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.writeXmlError(w, "SerializationException", err.Error(), http.StatusBadRequest)
		return
	}

	action := r.Form.Get("Action")
	if action == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	service := s.resolveQueryService(r.Header.Get("Authorization"), action)
	handler, ok := s.handlers[service]
	if !ok {
		s.writeXmlError(w, "UnknownService", "Service not implemented: "+service, http.StatusBadRequest)
		return
	}

	rc := common.FromContext(r.Context())
	response, err := handler.HandleQuery(action, r.Form, rc)
	if err != nil {
		s.writeXmlError(w, "InternalFailure", err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprint(w, response)
}

func (s *Server) resolveQueryService(auth, action string) string {
	if auth != "" {
		matches := servicePattern.FindStringSubmatch(auth)
		if len(matches) > 1 {
			scope := strings.ToLower(matches[1])
			if d := s.catalog.ByCredentialScope(scope); d != nil {
				return d.ExternalKey
			}
		}
	}

	// Fallback to infer from action if needed, but for now we expect auth for most SDKs.
	// In the Java version there's a large inference map.
	// For now, let's keep it simple and default to something if matched.
	return "sqs" // Default fallback for now
}

func (s *Server) writeXmlError(w http.ResponseWriter, code, message string, status int) {
	w.WriteHeader(status)
	xml := common.NewXmlBuilder().
		Start("ErrorResponse").
		Start("Error").
		Elem("Type", "Sender").
		Elem("Code", code).
		Elem("Message", message).
		End().
		Elem("RequestId", uuid.New().String()).
		End().
		Build()
	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprint(w, xml)
}

func (s *Server) handleJSONRequest(w http.ResponseWriter, r *http.Request) {
	target := r.Header.Get("X-Amz-Target")
	if target == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	match := s.catalog.MatchTarget(target)
	if match == nil {
		s.writeJSONError(w, "UnknownOperationException", "Unknown operation: "+target)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.writeJSONError(w, "SerializationException", err.Error())
		return
	}

	handler, ok := s.handlers[match.Descriptor.ExternalKey]
	if !ok {
		s.writeJSONError(w, "UnknownOperationException", "Service not implemented: "+match.Descriptor.ExternalKey)
		return
	}

	rc := common.FromContext(r.Context())
	response, err := handler.HandleJSON(match.Action, body, rc)
	if err != nil {
		// TODO: Map AwsException to JSON error
		s.writeJSONError(w, "InternalFailure", err.Error())
		return
	}

	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	json.NewEncoder(w).Encode(response)
}


func (s *Server) writeJSONError(w http.ResponseWriter, code, message string) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, `{"__type":"%s","message":"%s"}`, code, message)
}
