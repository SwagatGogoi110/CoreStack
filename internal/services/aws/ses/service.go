package ses

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/aws/ses/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type SesService struct {
	identityStore storage.Backend[string, *model.Identity]
	emailStore    storage.Backend[string, *model.SentEmail]
}

func NewSesService(factory *storage.Factory) (*SesService, error) {
	identityStore, _ := storage.CreateAccountAware[*model.Identity](factory, "ses", "ses-identities.json", "wal")
	emailStore, _ := storage.CreateAccountAware[*model.SentEmail](factory, "ses", "ses-emails.json", "wal")

	return &SesService{
		identityStore: identityStore,
		emailStore:    emailStore,
	}, nil
}

func (s *SesService) VerifyEmailIdentity(ctx context.Context, email string) (*model.Identity, error) {
	if _, ok, _ := s.identityStore.Get(ctx, email); ok {
		// Existing
	}

	identity := &model.Identity{
		Identity:           email,
		IdentityType:       "EmailAddress",
		VerificationStatus: "Success", // Auto-verify
		CreatedAt:          time.Now(),
	}

	if err := s.identityStore.Put(ctx, email, identity); err != nil {
		return nil, err
	}

	log.Printf("Verified SES email identity: %s", email)
	return identity, nil
}

func (s *SesService) ListIdentities(ctx context.Context) ([]string, error) {
	identities, err := s.identityStore.Scan(ctx, func(k string) bool { return true })
	if err != nil {
		return nil, err
	}
	var names []string
	for _, i := range identities {
		names = append(names, i.Identity)
	}
	return names, nil
}

func (s *SesService) SendEmail(ctx context.Context, source string, destinations []string, subject string) (string, error) {
	id := uuid.New().String()
	email := &model.SentEmail{
		MessageId:    id,
		Source:       source,
		Destinations: destinations,
		Subject:      subject,
		Timestamp:    time.Now(),
	}
	s.emailStore.Put(ctx, id, email)
	log.Printf("SES: Sent email from %s to %v: %s", source, destinations, subject)
	return id, nil
}
