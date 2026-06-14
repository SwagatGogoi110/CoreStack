package route53

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type Route53Handler struct {
	service *Route53Service
}

func NewRoute53Handler(service *Route53Service) *Route53Handler {
	return &Route53Handler{
		service: service,
	}
}

func (h *Route53Handler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("Route53 does not support standard JSON protocol dispatcher")
}

func (h *Route53Handler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("Route53 does not support Query protocol")
}

func (h *Route53Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/2013-04-01")
	
	if path == "/hostedzone" {
		if r.Method == "POST" {
			h.handleCreateHostedZone(w, r)
			return
		}
		if r.Method == "GET" {
			h.handleListHostedZones(w, r)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func (h *Route53Handler) handleCreateHostedZone(w http.ResponseWriter, r *http.Request) {
	// For simplicity, we create a zone with a mock name
	ctx := context.Background()
	zone, err := h.service.CreateHostedZone(ctx, "example.com.", "ref-123")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "<CreateHostedZoneResponse><HostedZone><Id>%s</Id><Name>%s</Name></HostedZone></CreateHostedZoneResponse>", zone.ID, zone.Name)
}

func (h *Route53Handler) handleListHostedZones(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	zones, _ := h.service.ListHostedZones(ctx)
	
	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprint(w, "<ListHostedZonesResponse><HostedZones>")
	for _, z := range zones {
		fmt.Fprintf(w, "<HostedZone><Id>%s</Id><Name>%s</Name></HostedZone>", z.ID, z.Name)
	}
	fmt.Fprint(w, "</HostedZones></ListHostedZonesResponse>")
}
