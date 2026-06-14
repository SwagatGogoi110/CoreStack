package privateca

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/privateca/model"
)

type CasHandler struct {
	service *CasService
}

func NewCasHandler(service *CasService) *CasHandler {
	return &CasHandler{
		service: service,
	}
}

func (h *CasHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/")
	ctx := context.Background()

	if strings.Contains(path, "/caPools") {
		parts := strings.Split(path, "/")
		
		if strings.HasSuffix(path, "/caPools") {
			if r.Method == "GET" {
				h.handleListPools(w, r, ctx, path)
				return
			}
			if r.Method == "POST" {
				h.handleCreatePool(w, r, ctx, path)
				return
			}
		}

		if strings.Contains(path, "/certificates") {
			idx := -1
			for i, p := range parts {
				if p == "caPools" {
					idx = i
					break
				}
			}
			poolName := strings.Join(parts[:idx+2], "/")
			if r.Method == "POST" {
				h.handleCreateCertificate(w, r, ctx, poolName)
				return
			}
		}

		if strings.Contains(path, "/certificateAuthorities") {
			idx := -1
			for i, p := range parts {
				if p == "caPools" {
					idx = i
					break
				}
			}
			poolName := strings.Join(parts[:idx+2], "/")
			if r.Method == "GET" {
				h.handleListCAs(w, r, ctx, poolName)
				return
			}
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "CAS path not implemented: %s"}}`, path)
}

func (h *CasHandler) handleListPools(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	items, _ := h.service.ListPools(ctx)
	res := model.CaPoolsList{
		CaPools: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *CasHandler) handleCreatePool(w http.ResponseWriter, r *http.Request, ctx context.Context, parent string) {
	var pool model.CaPool
	json.NewDecoder(r.Body).Decode(&pool)
	caPoolId := r.URL.Query().Get("caPoolId")
	name := parent + "/" + caPoolId
	created, _ := h.service.CreatePool(ctx, name, &pool)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *CasHandler) handleListCAs(w http.ResponseWriter, r *http.Request, ctx context.Context, poolName string) {
	items, _ := h.service.ListCAs(ctx, poolName)
	res := model.CertificateAuthoritiesList{
		CertificateAuthorities: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *CasHandler) handleCreateCertificate(w http.ResponseWriter, r *http.Request, ctx context.Context, poolName string) {
	var cert model.Certificate
	json.NewDecoder(r.Body).Decode(&cert)
	created, _ := h.service.CreateCertificate(ctx, poolName, &cert)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}
