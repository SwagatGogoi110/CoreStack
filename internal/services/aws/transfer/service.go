package transfer

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/aws/transfer/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

const (
	chars = "abcdefghijklmnopqrstuvwxyz0123456789"
)

type TransferService struct {
	serverStore storage.Backend[string, *model.Server]
}

func NewTransferService(factory *storage.Factory) (*TransferService, error) {
	serverStore, _ := storage.CreateAccountAware[*model.Server](factory, "transfer", "transfer-servers.json", "wal")

	return &TransferService{
		serverStore: serverStore,
	}, nil
}

func (s *TransferService) CreateServer(ctx context.Context, protocols []string) (*model.Server, error) {
	id := "s-" + s.randomID(17)
	server := &model.Server{
		ServerId:             id,
		Arn:                  fmt.Sprintf("arn:aws:transfer:us-east-1:000000000000:server/%s", id),
		State:                "ONLINE",
		Protocols:            protocols,
		EndpointType:         "PUBLIC",
		IdentityProviderType: "SERVICE_MANAGED",
		CreationTime:         time.Now(),
	}

	if err := s.serverStore.Put(ctx, id, server); err != nil {
		return nil, err
	}

	log.Printf("Created Transfer server: %s", id)
	return server, nil
}

func (s *TransferService) ListServers(ctx context.Context) ([]*model.Server, error) {
	return s.serverStore.Scan(ctx, func(k string) bool { return true })
}

func (s *TransferService) randomID(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
