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
	// First, check if this is a GCP request
	if s.isGcpRequest(r) {
		s.handleGCPRequest(w, r)
		return
	}

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

	// 2.1 Check for API Gateway V2
	if strings.HasPrefix(r.URL.Path, "/v2/apis") {
		// Route to apigatewayv2 handler
		body, _ := io.ReadAll(r.Body)
		rc := common.FromContext(r.Context())
		action := "GetApis"
		if r.Method == "POST" { action = "CreateApi" }
		res, _ := s.handlers["apigatewayv2"].HandleJSON(action, body, rc)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
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

	// 6. Check for Amazon MQ
	if strings.HasPrefix(r.URL.Path, "/v1/brokers") {
		s.mq.ServeHTTP(w, r)
		return
	}

	// 7. Generic REST-JSON routing for other services
	if r.Header.Get("Content-Type") == "application/json" || r.Header.Get("Content-Type") == "" {
		// Heuristics for common REST paths
		if strings.HasPrefix(r.URL.Path, "/applications") {
			s.handleGenericREST(w, r, "appconfig", "ListApplications")
			return
		}
		if strings.HasPrefix(r.URL.Path, "/v1/apis") {
			s.handleGenericREST(w, r, "appsync", "ListGraphqlApis")
			return
		}
		if strings.HasPrefix(r.URL.Path, "/backup/plans") {
			s.handleGenericREST(w, r, "backup", "ListBackupPlans")
			return
		}
		if strings.HasPrefix(r.URL.Path, "/clusters") {
			s.handleGenericREST(w, r, "eks", "ListClusters")
			return
		}
		if strings.HasPrefix(r.URL.Path, "/v1/clusters") {
			s.handleGenericREST(w, r, "msk", "ListClusters")
			return
		}
		if strings.HasPrefix(r.URL.Path, "/2021-01-01/domain") {
			s.handleGenericREST(w, r, "opensearch", "ListDomainNames")
			return
		}
		if strings.HasPrefix(r.URL.Path, "/v1/pipes") {
			s.handleGenericREST(w, r, "pipes", "ListPipes")
			return
		}
		if strings.HasPrefix(r.URL.Path, "/schedules") {
			s.handleGenericREST(w, r, "scheduler", "ListSchedules")
			return
		}
		if strings.HasPrefix(r.URL.Path, "/async-invoke") {
			s.handleGenericREST(w, r, "bedrock-runtime", "ListAsyncInvokes")
			return
		}
		if strings.HasPrefix(r.URL.Path, "/canaries") {
			s.handleGenericREST(w, r, "synthetics", "DescribeCanaries")
			return
		}
		if strings.HasPrefix(r.URL.Path, "/TraceSummaries") {
			s.handleGenericREST(w, r, "xray", "GetTraceSummaries")
			return
		}
	}

	// 8. Check for JSON protocols
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
			descriptors := s.catalog.ByCredentialScope(scope)
			if len(descriptors) == 1 {
				return descriptors[0].ExternalKey
			}
			if len(descriptors) > 1 {
				for _, d := range descriptors {
					for _, prefix := range d.TargetPrefixes {
						cleanPrefix := strings.TrimSuffix(prefix, ".")
						if strings.HasPrefix(action, cleanPrefix) {
							return d.ExternalKey
						}
					}
				}
				if scope == "rds" {
					if strings.Contains(action, "DBCluster") {
						return "docdb"
					}
					return "rds"
				}
				return descriptors[0].ExternalKey
			}
		}
	}
	return "sqs"
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

func (s *Server) handleGenericREST(w http.ResponseWriter, r *http.Request, service, action string) {
	handler, ok := s.handlers[service]
	if !ok {
		s.writeJSONError(w, "UnknownOperationException", "Service not implemented: "+service)
		return
	}

	body, _ := io.ReadAll(r.Body)
	rc := common.FromContext(r.Context())
	response, err := handler.HandleJSON(action, body, rc)
	if err != nil {
		s.writeJSONError(w, "InternalFailure", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// -- GCP Routing --

func (s *Server) isGcpRequest(r *http.Request) bool {
	if r.Header.Get("X-Goog-Api-Client") != "" {
		return true
	}
	if strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") && !strings.Contains(r.Header.Get("Authorization"), "AWS") {
		return true
	}
	path := r.URL.Path
	if strings.HasPrefix(path, "/v1/projects") || strings.HasPrefix(path, "/v2/projects") || 
	   strings.HasPrefix(path, "/v3/projects") ||
	   strings.HasPrefix(path, "/v1/locations") || strings.HasPrefix(path, "/v2/locations") ||
	   strings.HasPrefix(path, "/bigquery/v2/") || strings.HasPrefix(path, "/sql/v1/") ||
	   strings.HasPrefix(path, "/dns/v1/") || strings.HasPrefix(path, "/compute/v1/") ||
	   strings.HasPrefix(path, "/v1/apps/") || strings.HasPrefix(path, "/v2/entries") {
		return true
	}
	if strings.HasPrefix(path, "/storage/v1") || strings.HasPrefix(path, "/upload/storage") {
		return true
	}
	if strings.Contains(r.Host, "run.googleapis.com") {
		return true
	}
	return false
}

func (s *Server) handleGCPRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	var handler http.Handler

	if strings.HasPrefix(path, "/sql/v1/projects") {
		handler = s.gcpHandlers["cloudsql"]
	} else if strings.HasPrefix(path, "/dns/v1/projects") {
		handler = s.gcpHandlers["dns"]
	} else if strings.HasPrefix(path, "/v1/apps/") {
		handler = s.gcpHandlers["appengine"]
	} else if strings.HasPrefix(path, "/v2/entries") {
		handler = s.gcpHandlers["logging"]
	} else if strings.HasPrefix(path, "/compute/v1/projects") {
		if strings.Contains(path, "/securityPolicies") {
			handler = s.gcpHandlers["armor"]
		} else if strings.Contains(path, "/forwardingRules") || strings.Contains(path, "/targetHttpProxies") || strings.Contains(path, "/urlMaps") || strings.Contains(path, "/backendServices") {
			handler = s.gcpHandlers["loadbalancing"]
		} else {
			handler = s.gcpHandlers["compute"]
		}
	} else if strings.HasPrefix(path, "/v1/projects") || strings.HasPrefix(path, "/v2/projects") || strings.HasPrefix(path, "/v3/projects") {
		if strings.Contains(path, "/topics") || strings.Contains(path, "/subscriptions") {
			handler = s.gcpHandlers["pubsub"]
		} else if strings.Contains(path, "/secrets") {
			handler = s.gcpHandlers["secretmanager"]
		} else if strings.Contains(path, "/databases") {
			handler = s.gcpHandlers["firestore"]
		} else if strings.Contains(path, "/caPools") {
			handler = s.gcpHandlers["cas"]
		} else if strings.Contains(path, "/repositories") {
			handler = s.gcpHandlers["artifactregistry"]
		} else if strings.Contains(path, "/logs") || strings.Contains(path, "/entries") {
			handler = s.gcpHandlers["logging"]
		} else if strings.Contains(path, "/timeSeries") {
			handler = s.gcpHandlers["monitoring"]
		} else if strings.Contains(path, "/traces") {
			handler = s.gcpHandlers["trace"]
		} else if strings.Contains(path, "/instances") {
			if strings.HasPrefix(path, "/v2") {
				handler = s.gcpHandlers["bigtable"]
			} else {
				handler = s.gcpHandlers["spanner"]
			}
		} else if strings.Contains(path, "/serviceAccounts") {
			handler = s.gcpHandlers["iam"]
		} else if strings.Contains(path, "/operations") {
			handler = s.gcpHandlers["operations"]
		} else if strings.Contains(path, "/queues") {
			handler = s.gcpHandlers["tasks"]
		} else if strings.Contains(path, "/builds") || strings.Contains(path, "/triggers") {
			handler = s.gcpHandlers["cloudbuild"]
		} else if strings.Contains(path, "/workflows") {
			handler = s.gcpHandlers["workflows"]
		} else if strings.Contains(path, "/jobs") {
			handler = s.gcpHandlers["cloudscheduler"]
		} else if strings.Contains(path, "/functions") {
			handler = s.gcpHandlers["cloudfunctions"]
		} else if strings.Contains(path, "/services") {
			handler = s.gcpHandlers["cloudrun"]
		} else if strings.Contains(path, "/clusters") {
			if strings.Contains(path, "/locations/") {
				handler = s.gcpHandlers["gke"]
			} else {
				handler = s.gcpHandlers["kafka"]
			}
		} else if strings.Contains(path, "/bigquery/v2/") {
			handler = s.gcpHandlers["bigquery"]
		} else if strings.Contains(path, ":commit") || strings.Contains(path, ":lookup") || strings.Contains(path, ":runQuery") {
			handler = s.gcpHandlers["datastore"]
		}
	} else if strings.HasPrefix(path, "/bigquery/v2") {
		handler = s.gcpHandlers["bigquery"]
	} else if strings.HasPrefix(path, "/compute/v1") {
		handler = s.gcpHandlers["compute"]
	} else if strings.HasPrefix(path, "/storage/v1") || strings.HasPrefix(path, "/upload/storage") {
		handler = s.gcpHandlers["gcs"]
	}

	if handler != nil {
		handler.ServeHTTP(w, r)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "GCP service not implemented for path %s"}}`, path)
}
