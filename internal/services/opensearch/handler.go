package opensearch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type OpenSearchJsonHandler struct {
	service *OpenSearchService
}

func NewOpenSearchJsonHandler(service *OpenSearchService) *OpenSearchJsonHandler {
	return &OpenSearchJsonHandler{
		service: service,
	}
}

func (h *OpenSearchJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("OpenSearch does not support standard JSON protocol dispatcher")
}

func (h *OpenSearchJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("OpenSearch does not support Query protocol")
}

func (h *OpenSearchJsonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/2021-01-01")
	
	if path == "/opensearch/domain" && r.Method == "POST" {
		h.handleCreateDomain(w, r)
		return
	}
	if path == "/opensearch/domain" && r.Method == "GET" {
		h.handleListDomains(w, r)
		return
	}
	
	if strings.HasPrefix(path, "/opensearch/domain/") {
		name := strings.TrimPrefix(path, "/opensearch/domain/")
		if r.Method == "GET" {
			h.handleDescribeDomain(w, r, name)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func (h *OpenSearchJsonHandler) handleCreateDomain(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DomainName    string `json:"DomainName"`
		EngineVersion string `json:"EngineVersion"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	domain, err := h.service.CreateDomain(ctx, req.DomainName, req.EngineVersion)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"DomainStatus": domain})
}

func (h *OpenSearchJsonHandler) handleDescribeDomain(w http.ResponseWriter, r *http.Request, name string) {
	ctx := context.Background()
	domain, err := h.service.DescribeDomain(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"DomainStatus": domain})
}

func (h *OpenSearchJsonHandler) handleListDomains(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	names, _ := h.service.ListDomains(ctx)
	res := make([]map[string]string, 0)
	for _, n := range names {
		res = append(res, map[string]string{"DomainName": n})
	}
	json.NewEncoder(w).Encode(map[string]any{"DomainNames": res})
}
