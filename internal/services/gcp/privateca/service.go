package privateca

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/gcp/privateca/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CasService struct {
	poolStore storage.Backend[string, *model.CaPool]
	caStore   storage.Backend[string, *model.CertificateAuthority]
	certStore storage.Backend[string, *model.Certificate]
}

func NewCasService(factory *storage.Factory) (*CasService, error) {
	poolStore, _ := storage.CreateAccountAware[*model.CaPool](factory, "cas", "cas-pools.json", "wal")
	caStore, _ := storage.CreateAccountAware[*model.CertificateAuthority](factory, "cas", "cas-cas.json", "wal")
	certStore, _ := storage.CreateAccountAware[*model.Certificate](factory, "cas", "cas-certs.json", "wal")

	return &CasService{
		poolStore: poolStore,
		caStore:   caStore,
		certStore: certStore,
	}, nil
}

// Pools

func (s *CasService) CreatePool(ctx context.Context, name string, pool *model.CaPool) (*model.CaPool, error) {
	pool.Name = name
	s.poolStore.Put(ctx, name, pool)
	return pool, nil
}

func (s *CasService) ListPools(ctx context.Context) ([]*model.CaPool, error) {
	return s.poolStore.Scan(ctx, func(k string) bool { return true })
}

// Certificate Authorities

func (s *CasService) ListCAs(ctx context.Context, poolName string) ([]*model.CertificateAuthority, error) {
	return s.caStore.Scan(ctx, func(k string) bool { return true })
}

// Certificates

func (s *CasService) CreateCertificate(ctx context.Context, poolName string, cert *model.Certificate) (*model.Certificate, error) {
	id := uuid.New().String()
	cert.Name = poolName + "/certificates/" + id
	cert.CreateTime = time.Now()
	cert.PemCertificate = "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----"
	s.certStore.Put(ctx, cert.Name, cert)
	return cert, nil
}
