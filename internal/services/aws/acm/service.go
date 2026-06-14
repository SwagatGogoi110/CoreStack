package acm

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/aws/acm/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type AcmService struct {
	certStore storage.Backend[string, *model.Certificate]
}

func NewAcmService(factory *storage.Factory) (*AcmService, error) {
	certStore, _ := storage.CreateAccountAware[*model.Certificate](factory, "acm", "acm-certificates.json", "wal")

	return &AcmService{
		certStore: certStore,
	}, nil
}

func (s *AcmService) RequestCertificate(ctx context.Context, domainName string, sans []string) (*model.Certificate, error) {
	certID := uuid.New().String()
	arn := fmt.Sprintf("arn:aws:acm:us-east-1:000000000000:certificate/%s", certID)

	now := time.Now()
	issuedAt := now

	// Generate simple self-signed certificate
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: domainName,
		},
		NotBefore:             now,
		NotAfter:              now.AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              append([]string{domainName}, sans...),
	}

	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	
	certPEM := &bytes.Buffer{}
	pem.Encode(certPEM, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	
	keyPEM := &bytes.Buffer{}
	pem.Encode(keyPEM, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	cert := &model.Certificate{
		Arn:                     arn,
		DomainName:              domainName,
		SubjectAlternativeNames: sans,
		Status:                  model.StatusIssued,
		Type:                    model.TypeAmazonIssued,
		CreatedAt:               now,
		IssuedAt:                &issuedAt,
		NotBefore:               &template.NotBefore,
		NotAfter:                &template.NotAfter,
		KeyAlgorithm:            model.KeyAlgorithmRSA2048,
		CertificateBody:         certPEM.String(),
		PrivateKey:              keyPEM.String(),
	}

	if err := s.certStore.Put(ctx, certID, cert); err != nil {
		return nil, err
	}

	log.Printf("Requested ACM certificate: %s", arn)
	return cert, nil
}

func (s *AcmService) DescribeCertificate(ctx context.Context, arn string) (*model.Certificate, error) {
	id := s.extractID(arn)
	cert, ok, err := s.certStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("ResourceNotFoundException: Certificate not found")
	}
	return cert, nil
}

func (s *AcmService) ListCertificates(ctx context.Context) ([]*model.CertificateSummary, error) {
	certs, err := s.certStore.Scan(ctx, func(k string) bool { return true })
	if err != nil {
		return nil, err
	}
	var summaries []*model.CertificateSummary
	for _, c := range certs {
		summaries = append(summaries, &model.CertificateSummary{
			CertificateArn: c.Arn,
			DomainName:     c.DomainName,
		})
	}
	return summaries, nil
}

func (s *AcmService) extractID(arn string) string {
	// arn:aws:acm:region:account:certificate/id
	parts := bytes.Split([]byte(arn), []byte("/"))
	return string(parts[len(parts)-1])
}
