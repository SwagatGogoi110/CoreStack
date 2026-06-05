package cloudfront

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hectorvent/cloudstack/internal/core/common"
	"github.com/hectorvent/cloudstack/internal/services/cloudfront/model"
)

type CloudFrontHandler struct {
	service *CloudFrontService
}

func NewCloudFrontHandler(service *CloudFrontService) *CloudFrontHandler {
	return &CloudFrontHandler{
		service: service,
	}
}

func (h *CloudFrontHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("CloudFront does not support standard JSON protocol dispatcher")
}

func (h *CloudFrontHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("CloudFront does not support Query protocol")
}

func (h *CloudFrontHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/2020-05-31")
	
	if path == "/distribution" && r.Method == "POST" {
		h.handleCreateDistribution(w, r)
		return
	}
	if path == "/distribution" && r.Method == "GET" {
		h.handleListDistributions(w, r)
		return
	}
	
	if strings.HasPrefix(path, "/distribution/") {
		id := strings.TrimPrefix(path, "/distribution/")
		if r.Method == "GET" {
			h.handleGetDistribution(w, r, id)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func (h *CloudFrontHandler) handleCreateDistribution(w http.ResponseWriter, r *http.Request) {
	// For simplicity, we just create a minimal distribution
	ctx := context.Background()
	dist, err := h.service.CreateDistribution(ctx, &model.DistributionConfig{
		Enabled: true,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	// TODO: Return XML response
	fmt.Fprintf(w, "<Distribution><Id>%s</Id><Status>%s</Status><DomainName>%s</DomainName></Distribution>", dist.ID, dist.Status, dist.DomainName)
}

func (h *CloudFrontHandler) handleGetDistribution(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	dist, err := h.service.GetDistribution(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprintf(w, "<Distribution><Id>%s</Id><Status>%s</Status><DomainName>%s</DomainName></Distribution>", dist.ID, dist.Status, dist.DomainName)
}

func (h *CloudFrontHandler) handleListDistributions(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	dists, _ := h.service.ListDistributions(ctx)
	
	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprint(w, "<DistributionList>")
	for _, d := range dists {
		fmt.Fprintf(w, "<Items><DistributionSummary><Id>%s</Id><Status>%s</Status><DomainName>%s</DomainName></DistributionSummary></Items>", d.ID, d.Status, d.DomainName)
	}
	fmt.Fprint(w, "</DistributionList>")
}
