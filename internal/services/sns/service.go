package sns

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/sns/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type SnsService struct {
	topicStore storage.Backend[string, *model.Topic]
	subStore   storage.Backend[string, *model.Subscription]
	mu         sync.RWMutex
}

func NewSnsService(factory *storage.Factory) (*SnsService, error) {
	topicStore, _ := storage.CreateAccountAware[*model.Topic](factory, "sns", "sns-topics.json", "wal")
	subStore, _ := storage.CreateAccountAware[*model.Subscription](factory, "sns", "sns-subscriptions.json", "wal")

	return &SnsService{
		topicStore: topicStore,
		subStore:   subStore,
	}, nil
}

func (s *SnsService) CreateTopic(ctx context.Context, name string, attributes map[string]string) (*model.Topic, error) {
	arn := fmt.Sprintf("arn:aws:sns:us-east-1:000000000000:%s", name)
	if existing, ok, _ := s.topicStore.Get(ctx, arn); ok {
		return existing, nil
	}

	topic := &model.Topic{
		Name:       name,
		TopicArn:   arn,
		Attributes: attributes,
		CreatedAt:  time.Now(),
	}

	if err := s.topicStore.Put(ctx, arn, topic); err != nil {
		return nil, err
	}

	log.Printf("Created SNS topic: %s", arn)
	return topic, nil
}

func (s *SnsService) ListTopics(ctx context.Context) ([]*model.Topic, error) {
	return s.topicStore.Scan(ctx, func(k string) bool { return true })
}

func (s *SnsService) Subscribe(ctx context.Context, topicArn, protocol, endpoint string) (*model.Subscription, error) {
	if _, ok, _ := s.topicStore.Get(ctx, topicArn); !ok {
		return nil, fmt.Errorf("NotFound: Topic does not exist")
	}

	subArn := fmt.Sprintf("%s:%s", topicArn, uuid.New().String())
	sub := &model.Subscription{
		SubscriptionArn: subArn,
		TopicArn:        topicArn,
		Protocol:        protocol,
		Endpoint:        endpoint,
		Owner:           "000000000000",
		AccountID:       "000000000000",
	}

	if err := s.subStore.Put(ctx, subArn, sub); err != nil {
		return nil, err
	}

	log.Printf("Subscribed %s to %s", endpoint, topicArn)
	return sub, nil
}

func (s *SnsService) Publish(ctx context.Context, topicArn, message string) (string, error) {
	if _, ok, _ := s.topicStore.Get(ctx, topicArn); !ok {
		return "", fmt.Errorf("NotFound: Topic does not exist")
	}

	messageID := uuid.New().String()
	
	// TODO: Actually deliver to subscribers
	log.Printf("Published message %s to topic %s: %s", messageID, topicArn, message)
	
	return messageID, nil
}
